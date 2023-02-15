package store

import (
	"gorm.io/gorm"
)

type Store struct {
	DB *gorm.DB
}

type MessageStore struct {
}
