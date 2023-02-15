package api

import (
	"net/http"

	"github.com/RobleDev498/spaces/model"
)

// Spaces     		'api/spaces'
// Space         	'api/spaces/{space_id:[A-Za-z0-9@]+}'
// SpacesForUser 	'api/users/{user_id:[0-9]+}/spaces'
// SpaceMembers  	'api/spaces/{space_id:[0-9]+}/members'

func (api *Api) InitSpace() {
	api.BaseRoutes.SpacesForUser.Handle("", api.SessionHandler(getSpacesForUser)).Methods("GET")
	api.BaseRoutes.Space.Handle("", api.SessionHandler(getSpace)).Methods("GET")

}

func getSpacesForUser(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	if c.Session.UserID != c.Params.UserID {
		c.Err = c.NewPermissionError(model.PERMISSION_READ_USER_SPACES)
		return
	}

	spaces, err := c.App.GetSpacesForUser(c.Params.UserID)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(model.SpaceListToJson(spaces)))
}

func getSpace(c *Context, w http.ResponseWriter, r *http.Request) {

}
