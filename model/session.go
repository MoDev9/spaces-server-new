package model

import (
	"encoding/json"
	"io"
	"time"

	"gorm.io/datatypes"
)

const (
	SESSION_COOKIE_TOKEN = "SNAUTHTOKEN"
	SESSION_COOKIE_USER  = "SNUSERID"
	SESSION_COOKIE_CSRF  = "SNCSRF"

	SESSION_USER_AGENT = "User-agent"
	SESSION_IP         = "IP address"

	SESSION_WEB_EXPIRY_DAYS = 10
)

type Session struct {
	ID             string     `json:"id"`
	Token          string     `json:"token,omitempty"`
	CreatedAt      *EpochTime `json:"createdAt"`
	ExpiresAt      *EpochTime `json:"expiresAt,omitempty"`
	LastActivityAt *EpochTime `json:"lastActivityAt,omitempty"`
	UserID         string     `json:"userId,omitempty"`
	Props          datatypes.JSONMap
}

func (s *Session) ToJson() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func SessionFromJson(data io.Reader) *Session {
	var s *Session
	json.NewDecoder(data).Decode(&s)
	return s
}

func (s *Session) SetSessionExpiry(days int) {
	if s.CreatedAt == nil {
		exp := EpochTime(time.Now().AddDate(0, 0, days))
		s.ExpiresAt = &exp
	} else {
		exp := EpochTime(time.Time(*s.CreatedAt).AddDate(0, 0, days))
		s.ExpiresAt = &exp
	}
}

func (s *Session) PreUpdate() {
	t := EpochTime(time.Now())
	s.LastActivityAt = &t
}

func (s *Session) PreSave() {
	if s.ID == "" {
		s.ID = NewId()
	}

	if s.Token == "" {
		s.Token = NewId()
	}

	t := EpochTime(time.Now())
	s.CreatedAt = &t
	s.LastActivityAt = s.CreatedAt
}

func (s *Session) Sanitize() {
	s.Token = ""
}

func (s *Session) IsExpired() bool {

	/* if s.ExpiresAt <= 0 {
		return false
	} */

	now := time.Now()
	expires := time.Time(*s.ExpiresAt)

	return now.After(expires)
}

func (s *Session) AddProp(key string, value string) {
	if s.Props == nil {
		s.Props = make(map[string]interface{})
	}

	s.Props[key] = value
}

func (s *Session) GenerateCSRF() string {
	token := NewId()
	s.AddProp("csrf", token)
	return token
}

func (s *Session) GetCSRF() string {
	if s.Props == nil {
		return ""
	}

	return s.Props["csrf"].(string)
}

func SessionsToJson(o []*Session) string {
	b, err := json.Marshal(o)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func SessionsFromJson(data io.Reader) []*Session {
	var o []*Session
	json.NewDecoder(data).Decode(&o)
	return o
}
