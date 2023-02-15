package app

import (
	"fmt"
	"net/http"

	"github.com/RobleDev498/spaces/model"
)

type Context struct {
	App     *App
	Params  *Params
	Session model.Session
	Err     *model.AppErr

	ipAddress string
	path      string
	userAgent string
}

func (c *Context) RequireUserId() {
	if c.Err != nil {
		return
	}

	if c.Params.UserID == model.ME {
		c.Params.UserID = c.Session.UserID
	}

	if !model.IsValidId(c.Params.UserID) {
		c.SetInvalidUrlParam("user_id")
	}
}

func (c *Context) RequireStreamId() {
	if c.Err != nil {
		return
	}

	if !model.IsValidId(c.Params.StreamID) {
		c.SetInvalidUrlParam("stream_id")
	}
}

func (c *Context) RequireSpaceId() {
	if c.Err != nil {
		return
	}

	if !model.IsValidId(c.Params.SpaceID) {
		c.SetInvalidUrlParam("space_id")
	}
}

func (c *Context) SetInvalidParam(parameter string) {
	c.Err = NewInvalidParamError(parameter)
}

func (c *Context) SetInvalidUrlParam(param string) {
	c.Err = NewInvalidUrlParamError(param)
}

/* func (c *Context) SetPermissionError(param string) {

} */

func (c *Context) NewPermissionError(permissions ...string) *model.AppErr {
	var permissionsStr string
	for _, permission := range permissions {
		permissionsStr += permission
		permissionsStr += ","
	}
	return &model.AppErr{Msg: permissionsStr, StatusCode: http.StatusForbidden}
}

func NewInvalidParamError(parameter string) *model.AppErr {
	err := &model.AppErr{
		Msg:        "invalid body param error: " + parameter,
		StatusCode: http.StatusBadRequest,
	}
	return err
}

func (c *Context) RemoveSessionCookie(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     model.SESSION_COOKIE_TOKEN,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}

func (c *Context) SetSession(s *model.Session) {
	c.Session = *s
}

func NewInvalidUrlParamError(param string) *model.AppErr {
	err := &model.AppErr{
		Msg:        fmt.Sprintf("Invalid URL Param: %s", param),
		AppErrCode: model.ERROR_INVALID_URL_PARAM,
		StatusCode: http.StatusBadRequest,
	}

	return err
}
