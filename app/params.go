package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Params struct {
	UserID    string
	SpaceID   string
	StreamID  string
	MessageID string
}

func ParamsFromRequest(r *http.Request) *Params {
	params := &Params{}

	props := mux.Vars(r)
	//query := r.URL.Query()

	if val, ok := props["user_id"]; ok {
		params.UserID = val
	}

	if val, ok := props["space_id"]; ok {
		params.SpaceID = val
	}

	if val, ok := props["stream_id"]; ok {
		params.StreamID = val
	}

	if val, ok := props["message_id"]; ok {
		params.MessageID = val
	}
	return params
}
