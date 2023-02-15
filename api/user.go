package api

import (
	"net/http"

	"github.com/RobleDev498/spaces/model"
)

func (api *Api) InitUser() {
	api.BaseRoutes.Users.Handle("", api.Handler(createUser)).Methods("POST")
	api.BaseRoutes.User.Handle("", api.SessionHandler(getUser)).Methods("GET")

	//api.BaseRoutes.Users.Handle("", api.Handler(getUsers)).Methods("GET")
	api.BaseRoutes.Users.Handle("/login", api.Handler(login)).Methods("POST")
	api.BaseRoutes.Users.Handle("/logout", api.SessionHandler(logout)).Methods("POST")
	/* api.BaseRoutes.User.Handle("/friends", api.SessionHandler(getFriends)).Methods("GET")
	api.BaseRoutes.User.Handle("/friends", api.SessionHandler(addFriend)).Methods("POST") */
}

func createUser(c *Context, w http.ResponseWriter, r *http.Request) {
	user := model.UserFromJson(r.Body)

	ruser, err := c.App.CreateUserFromRegister(user)
	if err != nil {
		c.Err = err
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(ruser.ToJson()))
}

func addFriend(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	if c.Session.UserID != c.Params.UserID {
		c.Err = c.NewPermissionError(model.PERMISSION_READ_USER_SPACES)
		return
	}

	props := model.JsonToMap(r.Body)
	username, ok := props["username"].(string)
	if !ok {
		c.SetInvalidParam("username")
		return
	}

	user, err := c.App.GetUserByUsername(username)
	if err != nil {
		c.Err = err
		return
	}

	friend := &model.Friend{
		UserID:   c.Params.UserID,
		FriendID: user.ID,
	}

	ruser, err := c.App.AddFriend(friend)
	if err != nil {
		c.Err = err
		return
	}

	ruser.Status = "Pending"
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(ruser.ToJson()))
}

//Mock data
func getFriends(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	if c.Session.UserID != c.Params.UserID {
		c.Err = c.NewPermissionError(model.PERMISSION_READ_USER_SPACES)
		return
	}

	friends, err := c.App.GetFriends(c.Params.UserID)
	if err != nil {
		c.Err = err
		return
	}

	c.App.SanitizeFriends(friends)
	w.Write([]byte(model.UserListToJson(friends)))
}

func getUser(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	canView, err := c.App.UserCanBeViewed(c.Session.UserID, c.Params.UserID)
	if !canView || err != nil {
		c.Err = err
		return
	}

	user, err := c.App.GetUser(c.Params.UserID)
	if err != nil {
		c.Err = err
		return
	}

	if c.Session.ID == user.ID {
		user.Sanitize()
	}
	c.App.Store.UpdateSessionActivity(&c.Session)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(user.ToJson()))
}

func getUsers(c *Context, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(""))
}

func login(c *Context, w http.ResponseWriter, r *http.Request) {
	props := model.JsonToMap(r.Body)

	email := props["email"].(string)
	password := props["password"].(string)

	user, err := c.App.AuthenticateUser(email, password)
	if err != nil {
		c.Err = err
		return
	}

	err = c.App.Login(c, user, w, r)
	if err != nil {
		c.Err = err
		return
	}

	//if r.Header.Get(model.HEADER_REQUESTED_WITH) == model.HEADER_REQUESTED_WITH_XML {
	c.App.AttachSessionCookies(c, w, r)
	//}

	user.Sanitize()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(user.ToJson()))
}

func logout(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RemoveSessionCookie(w, r)

	if c.Session.ID != "" {
		if err := c.App.RevokeSessionById(c.Session.ID); err != nil {
			c.Err = err
			return
		}
	}

	m := make(map[string]string)
	m["STATUS"] = "OK"

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(model.MapToJson(m)))
}
