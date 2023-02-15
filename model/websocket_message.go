package model

import (
	"bytes"
	"encoding/json"
	"io"
)

const (
	WEBSOCKET_EVENT_TYPING                = "typing"
	WEBSOCKET_EVENT_POSTED                = "posted"
	WEBSOCKET_EVENT_MESSAGE_EDITED        = "message_edited"
	WEBSOCKET_EVENT_MESSAGE_DELETED       = "message_deleted"
	WEBSOCKET_EVENT_MESSAGE_UNREAD        = "message_unread"
	WEBSOCKET_EVENT_STREAM_CREATED        = "stream_created"
	WEBSOCKET_EVENT_STREAM_DELETED        = "stream_deleted"
	WEBSOCKET_EVENT_STREAM_UPDATED        = "stream_updated"
	WEBSOCKET_EVENT_STREAM_MEMBER_UPDATED = "stream_member_updated"
	WEBSOCKET_EVENT_NEW_USER              = "new_user"
	WEBSOCKET_EVENT_ADDED_TO_SPACE        = "added_to_space"
	WEBSOCKET_EVENT_LEAVE_SPACE           = "leave_space"
	WEBSOCKET_EVENT_UPDATE_SPACE          = "update_space"
	WEBSOCKET_EVENT_DELETE_SPACE          = "delete_space"
	WEBSOCKET_EVENT_RESTORE_SPACE         = "restore_space"
	WEBSOCKET_EVENT_UPDATE_SPACE_SCHEME   = "update_space_scheme"
	WEBSOCKET_EVENT_USER_ADDED            = "user_added"
	WEBSOCKET_EVENT_USER_UPDATED          = "user_updated"
	WEBSOCKET_EVENT_USER_ROLE_UPDATED     = "user_role_updated"
	WEBSOCKET_EVENT_USER_REMOVED          = "user_removed"
	WEBSOCKET_EVENT_STATUS_CHANGE         = "status_change"
	WEBSOCKET_EVENT_HELLO                 = "hello"
	WEBSOCKET_EVENT_REACTION_ADDED        = "reaction_added"
	WEBSOCKET_EVENT_REACTION_REMOVED      = "reaction_removed"
	WEBSOCKET_EVENT_RESPONSE              = "response"
	WEBSOCKET_EVENT_STREAM_VIEWED         = "stream_viewed"
	WEBSOCKET_EVENT_ROLE_UPDATED          = "role_updated"
)

const STATUS_FAIL = "FAIL"
const STATUS_OK = "OK"

type WebSocketMessage interface {
	ToJson() string
	EventType() string
	IsValid() bool
}

type WebSocketBroadcast struct {
	OmitUsers map[string]bool `json:"omitUsers,omitempty"` // broadcast is omitted for users listed here
	UserId    string          `json:"userId,omitempty"`    // broadcast only occurs for this user
	StreamID  string          `json:"streamId,omitempty"`  // broadcast only occurs for users in this stream
	SpaceID   string          `json:"spaceId,omitempty"`   // broadcast only occurs for users in this space
	Group     []string        `json:"group,omitempty"`     // broadcast nly occurs for users in this group
}

// webSocketEventJSON mirrors WebSocketEvent to make some of its unexported fields serializable
type webSocketEventJSON struct {
	Event     string                 `json:"event"`
	Data      map[string]interface{} `json:"data"`
	Broadcast *WebSocketBroadcast    `json:"broadcast"`
	Sequence  int64                  `json:"seq"`
}

type WebSocketEvent struct {
	Event     string
	Data      map[string]interface{}
	Broadcast *WebSocketBroadcast
	Sequence  int64
}

func NewWebSocketEvent(event, spaceId, streamId, userId string, omitUsers map[string]bool) *WebSocketEvent {
	return &WebSocketEvent{Event: event, Data: make(map[string]interface{}),
		Broadcast: &WebSocketBroadcast{SpaceID: spaceId, StreamID: streamId, UserId: userId, OmitUsers: omitUsers}}
}

func (ev *WebSocketEvent) Add(key string, value interface{}) {
	ev.Data[key] = value
}

func (ev *WebSocketEvent) Copy() *WebSocketEvent {
	copy := &WebSocketEvent{
		Event:     ev.Event,
		Data:      ev.Data,
		Broadcast: ev.Broadcast,
		Sequence:  ev.Sequence,
	}
	return copy
}

func (ev *WebSocketEvent) GetData() map[string]interface{} {
	return ev.Data
}

func (ev *WebSocketEvent) GetBroadcast() *WebSocketBroadcast {
	return ev.Broadcast
}

func (ev *WebSocketEvent) GetSequence() int64 {
	return ev.Sequence
}

func (ev *WebSocketEvent) SetEvent(event string) *WebSocketEvent {
	copy := ev.Copy()
	copy.Event = event
	return copy
}

func (ev *WebSocketEvent) SetData(data map[string]interface{}) *WebSocketEvent {
	copy := ev.Copy()
	copy.Data = data
	return copy
}

func (ev *WebSocketEvent) SetBroadcast(broadcast *WebSocketBroadcast) *WebSocketEvent {
	copy := ev.Copy()
	copy.Broadcast = broadcast
	return copy
}

func (ev *WebSocketEvent) SetSequence(seq int64) *WebSocketEvent {
	copy := ev.Copy()
	copy.Sequence = seq
	return copy
}

func (ev *WebSocketEvent) IsValid() bool {
	return ev.Event != ""
}

func (ev *WebSocketEvent) EventType() string {
	return ev.Event
}

func (ev *WebSocketEvent) ToJson() string {
	b, _ := json.Marshal(webSocketEventJSON{
		ev.Event,
		ev.Data,
		ev.Broadcast,
		ev.Sequence,
	})
	return string(b)
}

func WebSocketEventFromJson(data io.Reader) *WebSocketEvent {
	var ev WebSocketEvent
	var o webSocketEventJSON
	if err := json.NewDecoder(data).Decode(&o); err != nil {
		return nil
	}
	ev.Event = o.Event
	if u, ok := o.Data["user"]; ok {
		// We need to convert to and from JSON again
		// because the user is in the form of a map[string]interface{}.
		buf, err := json.Marshal(u)
		if err != nil {
			return nil
		}
		o.Data["user"] = UserFromJson(bytes.NewReader(buf))
	}
	ev.Data = o.Data
	ev.Broadcast = o.Broadcast
	ev.Sequence = o.Sequence
	return &ev
}

// WebSocketResponse represents a response received through the WebSocket
// for a request made to the server. This is available through the ResponseChannel
// channel in WebSocketClient.
type WebSocketResponse struct {
	Status   string                 `json:"status"`              // The status of the response. For example: OK, FAIL.
	SeqReply int64                  `json:"seq_reply,omitempty"` // A counter which is incremented for every response sent.
	Data     map[string]interface{} `json:"data,omitempty"`      // The data contained in the response.
	Error    *AppErr                `json:"error,omitempty"`     // A field that is set if any error has occurred.
}

func (m *WebSocketResponse) Add(key string, value interface{}) {
	m.Data[key] = value
}

func NewWebSocketResponse(status string, seqReply int64, data map[string]interface{}) *WebSocketResponse {
	return &WebSocketResponse{Status: status, SeqReply: seqReply, Data: data}
}

func NewWebSocketError(seqReply int64, err *AppErr) *WebSocketResponse {
	return &WebSocketResponse{Status: STATUS_FAIL, SeqReply: seqReply, Error: err}
}

func (m *WebSocketResponse) IsValid() bool {
	return m.Status != ""
}

func (m *WebSocketResponse) EventType() string {
	return WEBSOCKET_EVENT_RESPONSE
}

func (m *WebSocketResponse) ToJson() string {
	b, _ := json.Marshal(m)
	return string(b)
}

func WebSocketResponseFromJson(data io.Reader) *WebSocketResponse {
	var o *WebSocketResponse
	json.NewDecoder(data).Decode(&o)
	return o
}
