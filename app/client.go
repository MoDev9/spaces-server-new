package app

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/RobleDev498/spaces/model"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// The websocket connection.
	App  *App
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan model.WebSocketMessage

	Session      *model.Session
	UserId       string
	ConnectionID string
	Sequence     int64
}

func NewClient(conn *websocket.Conn, session *model.Session, seq int64, a *App) *Client {
	client := &Client{
		conn:     conn,
		App:      a,
		send:     make(chan model.WebSocketMessage),
		Session:  session,
		Sequence: seq,
	}

	client.UserId = client.Session.UserID
	client.ConnectionID = model.NewId()
	return client
}

func (c *Client) createHelloMessage() *model.WebSocketEvent {
	msg := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_HELLO, "", "", c.UserId, nil)
	msg.Add("message", "Hell0 World")
	msg.Add("connectionId", c.ConnectionID)
	return msg
}

// shouldSendEvent returns whether the message should be sent or not.
func (c *Client) shouldSendEvent(msg *model.WebSocketEvent) bool {
	// If the event is destined to a specific user
	if msg.GetBroadcast().UserId != "" {
		return c.UserId == msg.GetBroadcast().UserId
	}

	if msg.GetBroadcast().StreamID != "" {
		result, err := c.App.GetStreamMembers(msg.GetBroadcast().StreamID)
		if err != nil {
			for _, v := range result {
				if c.UserId == v.ID {
					return true
				}
			}
		}
	}

	return false
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		//c.hub.unregister <- c
		c.App.HubUnRegister(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		data := make(map[string]interface{})
		data["message"] = string(message)
		ev := &model.WebSocketEvent{
			Data:      data,
			Broadcast: &model.WebSocketBroadcast{},
		}
		c.App.Publish(ev)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			msg := message.(*model.WebSocketEvent)
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			//w.Write([]byte(msg.Data["message"].(string)))
			w.Write([]byte(msg.ToJson()))

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				messsage := <-c.send
				msg = messsage.(*model.WebSocketEvent)
				// w.Write([]byte(msg.Data["message"].(string)))
				w.Write([]byte(msg.ToJson()))
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(c *Context, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("active goroutines : connWebSocket %d\n", runtime.NumGoroutine())
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     c.App.OriginChecker(),
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(conn, &c.Session, 0, c.App)
	if c.Session.UserID != "" {
		c.App.HubRegister(client)
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
