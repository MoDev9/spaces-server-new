package api

import (
	"net/http"

	"github.com/RobleDev498/spaces/model"
)

func (api *Api) InitMessage() {
	api.BaseRoutes.Messages.Handle("", api.SessionHandler(createMessage)).Methods("POST")

	api.BaseRoutes.MessagesForStream.Handle("", api.SessionHandler(getMessages)).Methods("GET")
}

func createMessage(c *Context, w http.ResponseWriter, r *http.Request) {
	msg := model.MessageFromJson(r.Body)
	if msg == nil {
		c.SetInvalidParam("post")
		return
	}

	msg.AuthorID = c.Session.UserID
	rMessage, err := c.App.CreateMessage(c, msg)
	if err != nil {
		c.Err = err
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(rMessage.ToJson()))
}

func getMessages(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireStreamId()
	if c.Err != nil {
		return
	}

	messages, err := c.App.GetMessagesForStream(c, c.Session.UserID, c.Params.StreamID)

	if err != nil {
		c.Err = nil
		return
	}

	w.Write([]byte(model.MessageListToJson(messages)))
}
