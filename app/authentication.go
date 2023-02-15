package app

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/RobleDev498/spaces/model"
)

const (
	HEADER_AUTH  = "Authorization"
	TOKEN_BEARER = "Bearer"
	HEADER_TOKEN = "Token"

	TOKEN_LOCATION_HEADER = 1
	TOKEN_LOCATION_COOKIE = 2
)

var restrictedNames = []string{"everyone"}

func IsPasswordValid(password string) *model.AppErr {
	/* if len(password) <= model.MINIMUM_LENGTH_PASSWORD {
		return &model.AppErr{
			Msg:        "Password length must be equal to or greater than 8 ",
			AppErrCode: model.ERROR_PASSWORD_LENGTH,
			StatusCode: http.StatusUnprocessableEntity,
		}
	} */
	return nil
}

func IsUsernameValid(username string) bool {
	l := len(username)
	if l < model.MINIMUM_LENGTH_USERNAME {
		return false
	} else if l > model.MAXIMUM_LENGTH_USERNAME {
		return false
	}

	return true
}

func (a *App) ValidUsername(userName string) *model.AppErr {
	username := strings.ToLower(userName)
	if !IsUsernameValid(userName) {
		return &model.AppErr{
			Msg:        "Username length must be between 2 and 32 characters",
			AppErrCode: model.ERROR_USERNAME_INVALID,
			StatusCode: http.StatusUnprocessableEntity,
		}
	}

	appErr := &model.AppErr{
		Msg:        fmt.Sprintf("Username %s taken", userName),
		AppErrCode: model.ERROR_USERNAME_TAKEN,
		StatusCode: http.StatusOK,
	}

	for _, s := range restrictedNames {
		if username == s {
			return appErr
		}
	}

	return nil
}

func ParseToken(r *http.Request) (string, int) {
	authHeader := r.Header.Get(HEADER_AUTH)

	if authHeader != "" && authHeader[0:6] == TOKEN_BEARER {
		return authHeader[7:], TOKEN_LOCATION_HEADER
	}

	if cookie, err := r.Cookie(model.SESSION_COOKIE_TOKEN); err == nil {
		return cookie.Value, TOKEN_LOCATION_COOKIE
	}

	return "", -1
}
