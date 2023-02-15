package store

import (
	"log"
	"time"

	"github.com/RobleDev498/spaces/model"
)

func (s *Store) CreateSession(session *model.Session) (*model.Session, error) {
	session.Token = ""
	session.PreSave()
	result := s.DB.Create(&session)
	return session, result.Error
}

func (s *Store) GetSession(token string) (*model.Session, error) {
	var session *model.Session
	result := s.DB.Table("sessions").Where("token = ?", token).First(&session)
	return session, result.Error
}

func (s *Store) RemoveSession(session *model.Session) error {
	result := s.DB.Delete(&session)
	return result.Error
}

func (s *Store) GetSessionById(id string) (*model.Session, error) {
	var session *model.Session
	result := s.DB.First(&session, "id = ?", id)
	return session, result.Error
}

func (s *Store) UpdateSessionActivity(session *model.Session) {
	result := s.DB.Model(session).Update("last_activity_at", time.Now())
	if err := result.Error; err != nil {
		log.Println(err)
	}
}
