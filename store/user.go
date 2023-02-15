package store

import (
	"strings"

	"github.com/RobleDev498/spaces/model"
)

func (s *Store) UpdateUser(u *model.User) (*model.User, error) {
	result := s.DB.Save(u)
	return u, result.Error
}

func (s *Store) CreateUser(u *model.User) (*model.User, error) {
	u.PreSave()
	result := s.DB.Create(&u)
	//u.Sanitize()
	return u, result.Error
}

func (s *Store) GetUser(id string) (*model.User, error) {
	var u *model.User
	result := s.DB.First(&u, "id = ?", id)
	return u, result.Error
}

func (s *Store) GetAllUsers() ([]*model.User, error) {
	var u []*model.User
	result := s.DB.Find(u)
	return u, result.Error
}

func (s *Store) GetUserByUsername(Username string) (*model.User, error) {
	username := strings.ToLower(Username)
	var u *model.User
	result := s.DB.Table("users").Where("LOWER(username) = ?", username).First(&u)
	return u, result.Error
}

func (s *Store) GetUserByEmail(Email string) (*model.User, error) {
	email := strings.ToLower(Email)
	var u *model.User
	result := s.DB.Table("users").Where("LOWER(email) = ?", email).First(&u)
	return u, result.Error
}
