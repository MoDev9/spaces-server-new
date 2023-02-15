package model

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	MINIMUM_LENGTH_USERNAME = 2
	MAXIMUM_LENGTH_USERNAME = 32

	ME = "me"

	STATUS_FRIEND_ACCEPTED = 11
	STATUS_FRIEND_REQUEST  = 12
)

type EpochTime time.Time

type User struct {
	ID        string     `json:"id"`
	CreatedAt *EpochTime `json:"createdAt"`
	UpdatedAt *EpochTime `json:"updatedAt"`
	Email     string     `json:"email"`
	FirstName string     `json:"firstName,omitempty"`
	LastName  string     `json:"lastName,omitempty"`
	Username  string     `json:"username"`
	Password  string     `json:"password,omitempty"`
	Status    string     `json:"status"`
	//Space     Space `gorm:"foreignKey:HomeID"`
}

func (u *User) PreSave() {
	if u.ID == "" {
		u.ID = NewId()
	}

	t := EpochTime(time.Now())
	u.CreatedAt = &t

	if u.Password != "" {
		u.Password = HashPassword(u.Password)
	}
}

func (u *User) PreUpdate() {
	t := EpochTime(time.Now())
	u.UpdatedAt = &t
}

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Panic(err)
	}

	return string(hash)
}

func CheckPasswordHash(hash, password string) bool {
	if password == "" || hash == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (et *EpochTime) UnmarshalJSON(data []byte) error {
	t := strings.Trim(string(data), `"`) // Remove quote marks from around the JSON string
	sec, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return err
	}
	epochTime := time.Unix(sec, 0)
	*et = EpochTime(epochTime)
	return nil
}

func (et EpochTime) MarshalJSON() ([]byte, error) {
	timestamp := fmt.Sprintf("\"%d\"", time.Time(et).Unix())
	return []byte(timestamp), nil
}

func (u *User) Sanitize() {
	u.Password = ""
	//u.HomeID = ""
}

func (u *User) SanitizeFriend() {
	u.Sanitize()

	u.Email = ""
	u.CreatedAt = nil
	u.UpdatedAt = nil
	u.FirstName = ""
	u.LastName = ""
}

func UserListToJson(u []*User) string {
	b, _ := json.Marshal(u)
	return string(b)
}

func UserListFromJson(data io.Reader) []*User {
	var users []*User
	json.NewDecoder(data).Decode(&users)
	return users
}

func UserFromJson(data io.Reader) *User {
	var user *User
	json.NewDecoder(data).Decode(&user)
	return user
}

//Convert a User to JSON string
func (u *User) ToJson() string {
	b, _ := json.Marshal(u)
	return string(b)
}
