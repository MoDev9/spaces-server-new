package api

import (
	"net/http"

	"github.com/RobleDev498/spaces/app"
)

type Context = app.Context

func (api *Api) Handler(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	handler := &app.Handler{
		App:            api.app,
		HandleFunc:     h,
		RequireSession: false,
	}

	return handler
}

func (api *Api) SessionHandler(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	handler := &app.Handler{
		App:            api.app,
		HandleFunc:     h,
		RequireSession: true,
	}

	return handler
}
