package app

import (
	"github.com/RobleDev498/spaces/model"
)

// Hub maintains all websocket connections.
type Hub struct {
	app        *App
	broadcast  chan *model.WebSocketEvent // Inbound messages from the clients.
	register   chan *Client               // Register requests from the clients.
	unregister chan *Client               // Unregister requests from the clients.
	//Map UserId to Websocket Connection
	clients map[string]*Client
}

func (a *App) NewHub() *Hub {
	return &Hub{
		app:        a,
		broadcast:  make(chan *model.WebSocketEvent),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

func (s *Server) GetHub() *Hub {
	return s.hub
}

func (a *App) GetHub() *Hub {
	return a.Srv().GetHub()
}

func (a *App) HubRegister(c *Client) {
	hub := a.GetHub()
	if hub != nil {
		hub.Register(c)
	}
}

func (a *App) HubUnRegister(c *Client) {
	hub := a.GetHub()
	if hub != nil {
		hub.UnRegister(c)
	}
}

func (hub *Hub) Register(c *Client) {
	hub.register <- c
}

func (hub *Hub) UnRegister(c *Client) {
	hub.unregister <- c
}

func (h *Hub) Broadcast(message *model.WebSocketEvent) {
	if h != nil && message != nil {
		select {
		case h.broadcast <- message:
		default:
		}
	}
}

func (a *App) Publish(ev *model.WebSocketEvent) {
	a.Srv().hub.Broadcast(ev)
}

func (a *App) StartHub() {
	hub := a.NewHub()
	a.srv.hub = hub
	go hub.Start()
}

func (h *Hub) Has(c *Client) bool {
	_, ok := h.clients[c.UserId]
	return ok
}

func (h *Hub) Add(c *Client) {
	h.clients[c.UserId] = c
}

func (h *Hub) Remove(c *Client) {
	delete(h.clients, c.UserId)
}

func (hub *Hub) Start() {

	for {
		select {
		case client := <-hub.register:
			if hub.Has(client) {
				hub.Remove(client)
			}
			hub.Add(client)
			client.send <- client.createHelloMessage()
			//}
		case client := <-hub.unregister:
			hub.Remove(client)
			close(client.send)
		case message := <-hub.broadcast:
			broadcast := func(client *Client) {
				if !hub.Has(client) {
					return
				}

				if client.shouldSendEvent(message) {
					select {
					case client.send <- message:
					default:
						close(client.send)
						hub.Remove(client)
					}
				}
			}

			if message.GetBroadcast().UserId != "" {
				client, ok := hub.clients[(message.GetBroadcast().UserId)]
				if ok {
					broadcast(client)
				}

				continue
			} else if message.GetBroadcast().StreamID != "" {
				members, err := hub.app.GetStreamMembers(message.GetBroadcast().StreamID)
				if err != nil {
					for _, member := range members {
						client, ok := hub.clients[member.ID]
						if ok {
							broadcast(client)
						}
					}
				}
			} else {
				for _, client := range hub.clients {
					broadcast(client)
				}
			}
		}
	}
}
