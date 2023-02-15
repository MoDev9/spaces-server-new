package model

import (
	"encoding/json"
	"io"
	"time"
)

const (
	GUILD_TEXT  int = 1
	DM          int = 2
	GUILD_VOICE int = 3
	GROUP_DM    int = 4

	STREAM_TOPIC_LENGTH = 1024
	NAME_MINIMUM        = 2
	NAME_MAXIMUM        = 100
)

type RoomMember struct {
	UserID   string `gorm:"user_id" json:"userId"`
	StreamID string `gorm:"stream_id" json:"streamId"`
}

type Stream struct {
	ID          string     `json:"id"`
	CreatedAt   *EpochTime `json:"createdAt"`
	UpdatedAt   *EpochTime `json:"updatedAt,omitempty"`
	Visibility  int        `json:"visibility,omitempty"`
	Type        int        `json:"type"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	SpaceID     string     `json:"spaceId,omitempty"`
	OwnerID     string     `json:"ownerId,omitempty"`
	Members     []*User    `json:"members,omitempty" gorm:"-"`
	MemberCount int        `json:"memberCount,omitempty"`
	//Space       Space
}

func (s *Stream) PreSave() {
	if s.ID == "" {
		s.ID = NewId()
	}

	t := EpochTime(time.Now())
	s.CreatedAt = &t
}

func (s *Stream) PreUpdate() {
	t := EpochTime(time.Now())
	s.UpdatedAt = &t
}

func RoomMemberListToJson(r []*RoomMember) string {
	b, _ := json.Marshal(r)
	return string(b)
}

func RoomMemberListFromJson(data io.Reader) []*RoomMember {
	var roomMembers []*RoomMember
	json.NewDecoder(data).Decode(&roomMembers)
	return roomMembers
}

func StreamFromJson(data io.Reader) *Stream {
	var stream *Stream
	json.NewDecoder(data).Decode(&stream)
	return stream
}

func StreamListFromJson(data io.Reader) []*Stream {
	var streams []*Stream
	json.NewDecoder(data).Decode(&streams)
	return streams
}

func StreamListToJson(s []*Stream) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s *Stream) ToJson() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (r *RoomMember) ToJson() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (s *Stream) Sanitize() {

}
