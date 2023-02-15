package model

import (
	"encoding/json"
	"io"
)

type Space struct {
	ID          string     `json:"id"`
	CreatedAt   *EpochTime `json:"createdAt"`
	UpdatedAt   *EpochTime `json:"updatedAt,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Link        string     `json:"link,omitempty"`
	OwnerID     string     `json:"ownerId,omitempty"`
	//Users       []User
}

func (space *Space) PreSave() {
	if space.ID == "" {
		space.ID = NewId()
	}
}

func (space *Space) Sanitize() {
	space.OwnerID = ""
}

func SpaceFromJson(data io.Reader) *Space {
	var space *Space
	json.NewDecoder(data).Decode(&space)
	return space
}

func SpaceListToJson(s []*Space) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s *Space) ToJson() string {
	b, _ := json.Marshal(s)
	return string(b)
}
