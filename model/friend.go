package model

import (
	"encoding/json"
	"io"
	"time"
)

type Friend struct {
	UserID    string     `gorm:"column:user_id" json:"userId"`
	FriendID  string     `gorm:"column:friend_id" json:"friendId"`
	CreatedAt *EpochTime `gorm:"column:created_at" json:"createdAt"`
	Status    int        `gorm:"column:status" json:"status"`
}

func NewFriend(userId, friendId string) *Friend {
	return &Friend{
		UserID:   userId,
		FriendID: friendId,
	}
}

func (friend *Friend) ToJson() string {
	b, _ := json.Marshal(friend)
	return string(b)
}

func FriendFromJson(data io.Reader) *Friend {
	var friend *Friend
	json.NewDecoder(data).Decode(&friend)
	return friend
}

func (f *Friend) PreSave() {
	t := EpochTime(time.Now())
	f.CreatedAt = &t

	if f.UserID > f.FriendID {
		userId := f.UserID
		f.UserID = f.FriendID
		f.FriendID = userId
	}

	f.Status = STATUS_FRIEND_REQUEST
}
