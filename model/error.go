package model

import (
	"encoding/json"
	"log"
)

const (
	ERROR_USERNAME_TAKEN          = 101
	ERROR_EMAIL_EXISTS            = 102
	ERROR_USERNAME_INVALID        = 103
	ERROR_PASSWORD_LENGTH         = 104
	ERROR_INVALID_LOGIN_CRED      = 105
	ERROR_INVALID_EXPIRED_SESSION = 106

	ERROR_RECORD_NOT_FOUND = 41
	ERROR_INTERNAL_ERROR   = 500

	ERROR_INVALID_URL_PARAM = 251

	ERROR_PG_UNIQUE           = "23505"
	ERROR_CONSTRAINT_USERNAME = "users_username_key"
	ERROR_CONSTRAINT_EMAIL    = "users_email_key"
)

type AppErr struct {
	Msg        string
	AppErrCode int
	StatusCode int
}

func (appErr *AppErr) ToJson() string {
	b, err := json.Marshal(appErr)
	if err != nil {
		log.Panic(err)
	}

	return string(b)
}
