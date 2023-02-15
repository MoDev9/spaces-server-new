package model

import (
	"encoding/json"
	"io"
	"time"
)

type Message struct {
	ID        string     `json:"id"`
	CreatedAt *EpochTime `json:"createdAt"`
	UpdatedAt *EpochTime `json:"updatedAt,omitempty"`
	Content   string     `json:"content,omitempty"`
	AuthorID  string     `json:"authorId"`
	StreamID  string     `json:"streamId"`
	//Author    User
	//Stream    Stream
}

func (m *Message) PreSave() {
	if m.ID == "" {
		m.ID = NewId()
	}

	t := EpochTime(time.Now())
	m.CreatedAt = &t
}

func (m *Message) ToJson() string {
	b, _ := json.Marshal(m)
	return string(b)
}

func MessageListToJson(m []*Message) string {
	b, _ := json.Marshal(m)
	return string(b)
}

func MessageFromJson(data io.Reader) *Message {
	var msg *Message
	json.NewDecoder(data).Decode(&msg)
	return msg
}
