package api

import (
	"net/http"

	"github.com/RobleDev498/spaces/model"
)

func (api *Api) InitStream() {
	api.BaseRoutes.RoomsForUser.Handle("", api.SessionHandler(getRooms)).Methods("GET")
	api.BaseRoutes.StreamMembers.Handle("", api.SessionHandler(getStreamMembers)).Methods("GET")
	api.BaseRoutes.StreamMembers.Handle("", api.SessionHandler(addUserToRoom)).Methods("POST")
	//api.BaseRoutes.RoomsForUser.Handle("", api.SessionHandler(getRooms)).Methods("GET")

	api.BaseRoutes.Streams.Handle("/default", api.SessionHandler(getDefaultStream)).Methods("GET")

	//api.BaseRoutes.Streams.Handle("", api.Handler(getRooms)).Methods("GET")
	api.BaseRoutes.Streams.Handle("", api.SessionHandler(createStream)).Methods("POST")
	api.BaseRoutes.Stream.Handle("", api.SessionHandler(getStream)).Methods("GET")
}

func createStream(c *Context, w http.ResponseWriter, r *http.Request) {
	stream := model.StreamFromJson(r.Body)

	rstream, err := c.App.CreateStream(c.Session.UserID, stream)
	if err != nil {
		c.Err = err
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(rstream.ToJson()))
}

func getStream(c *Context, w http.ResponseWriter, r *http.Request) {

}

func getDefaultStream(c *Context, w http.ResponseWriter, r *http.Request) {
	streamId, err := c.App.GetDefaultStream(c.Session.UserID)
	if err != nil {
		c.Err = err
		return
	}

	jsonMap := make(map[string]string)
	jsonMap["id"] = streamId

	w.Write([]byte(model.MapToJson(jsonMap)))
}

func getStreamMembers(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireStreamId()
	if c.Err != nil {
		return
	}

	users, err := c.App.GetStreamMembers(c.Params.StreamID)
	if err != nil {
		c.Err = err
		return
	}

	c.App.SanitizeFriends(users)
	w.Write([]byte(model.UserListToJson(users)))
}

func addUserToRoom(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireStreamId()
	if c.Err != nil {
		return
	}

	props := model.JsonToMap(r.Body)
	username, ok := props["username"].(string)
	if !ok {
		c.SetInvalidParam("username")
		return
	}

	ruser, err := c.App.GetUserByUsername(username)
	if err != nil {
		c.Err = err
		return
	}

	rm := &model.RoomMember{
		UserID:   ruser.ID,
		StreamID: c.Params.StreamID,
	}

	_, err = c.App.AddUserToRoom(c.Session.UserID, rm, ruser)
	if err != nil {
		c.Err = err
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(ruser.ToJson()))
}

func getRooms(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	if c.Session.UserID != c.Params.UserID {
		c.Err = c.NewPermissionError(model.PERMISSION_READ_USER_SPACES)
		return
	}

	rooms, err := c.App.GetRooms(c.Params.UserID)
	if err != nil {
		c.Err = err
		return
	}

	c.App.SanitizeRooms(c.Session, rooms)
	w.Write([]byte(model.StreamListToJson(rooms)))
}
