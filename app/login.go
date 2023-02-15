package app

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/RobleDev498/spaces/model"
	"gorm.io/gorm"
)

func (app *App) AuthenticateUser(email, password string) (*model.User, *model.AppErr) {
	email = strings.ToLower(strings.TrimSpace(email))

	/* if err := IsPasswordValid(password); err != nil {
		return nil, err
	} */
	invalidLoginErr := &model.AppErr{
		Msg:        "Enter a valid email or username and/or password",
		AppErrCode: model.ERROR_INVALID_LOGIN_CRED,
		StatusCode: http.StatusUnauthorized,
	}

	user, err := app.Store.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, invalidLoginErr
		}

		return nil, &model.AppErr{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
			AppErrCode: model.ERROR_INTERNAL_ERROR,
		}
	}

	if !model.CheckPasswordHash(user.Password, password) {
		return nil, invalidLoginErr
	}

	return user, nil
}

func (app *App) Login(c *Context, user *model.User, w http.ResponseWriter, r *http.Request) *model.AppErr {
	session := &model.Session{
		UserID: user.ID,
		Props: map[string]interface{}{
			model.SESSION_USER_AGENT: r.UserAgent(),
			model.SESSION_IP:         r.RemoteAddr,
		},
	}
	session.GenerateCSRF()
	session.SetSessionExpiry(model.SESSION_WEB_EXPIRY_DAYS)

	var err error
	if session, err = app.Store.CreateSession(session); err != nil {
		return &model.AppErr{
			StatusCode: http.StatusInternalServerError,
		}
	}

	w.Header().Add(HEADER_TOKEN, session.Token)
	c.SetSession(session)

	user.Status = "Online"
	c.App.Store.UpdateUser(user)
	return nil
}

func (a *App) AttachSessionCookies(c *Context, w http.ResponseWriter, r *http.Request) {
	maxAge := model.SESSION_WEB_EXPIRY_DAYS * 24 * 60 * 60
	expiresAt := time.Now().AddDate(0, 0, model.SESSION_WEB_EXPIRY_DAYS)

	sessionCookie := &http.Cookie{
		Name:     model.SESSION_COOKIE_TOKEN,
		Value:    c.Session.Token,
		MaxAge:   maxAge,
		Expires:  expiresAt,
		HttpOnly: true,
		Path:     "/",
	}

	csrfCookie := &http.Cookie{
		Name:    model.SESSION_COOKIE_CSRF,
		Value:   c.Session.GetCSRF(),
		MaxAge:  maxAge,
		Expires: expiresAt,
		Path:    "/",
	}

	userCookie := &http.Cookie{
		Name:    model.SESSION_COOKIE_USER,
		Value:   c.Session.UserID,
		MaxAge:  maxAge,
		Expires: expiresAt,
		Path:    "/",
	}

	http.SetCookie(w, sessionCookie)
	http.SetCookie(w, csrfCookie)
	http.SetCookie(w, userCookie)
}
