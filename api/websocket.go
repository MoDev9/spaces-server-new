package api

import (
	"github.com/RobleDev498/spaces/app"
)

const (
	connectionIDParam   = "connection_id"
	sequenceNumberParam = "sequence_number"
)

func (api *Api) InitWebSocket() {
	api.BaseRoutes.ApiRoot.Handle("/{websocket:websocket(?:\\/)?}", api.SessionHandler(app.ServeWs))
}
