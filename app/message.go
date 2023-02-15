package app

import (
	"log"
	"net/http"

	"github.com/RobleDev498/spaces/model"
)

func (a *App) CreateMessage(c *Context, msg *model.Message) (*model.Message, *model.AppErr) {
	stream, err := a.GetStream(msg.StreamID)
	if err != nil {
		return nil, err
	}

	user, err := a.GetUser(msg.AuthorID)
	if err != nil {
		return nil, err
	}

	rm, err2 := a.createMessage(msg)
	if err2 != nil {
		return nil, &model.AppErr{
			Msg:        err2.Error(),
			StatusCode: http.StatusInternalServerError,
			AppErrCode: model.ERROR_INTERNAL_ERROR,
		}
	}

	if err2 := a.handleMessageEvents(c, rm, user, stream); err2 != nil {
		log.Println(err2)
	}

	return rm, nil
}

func (a *App) handleMessageEvents(c *Context, msg *model.Message, user *model.User, stream *model.Stream) error {
	var space *model.Space
	var err error
	if stream.SpaceID != "" {
		space, err = a.getSpace(stream.SpaceID)
		if err != nil {
			return err
		}
	}

	return a.SendNotifications(msg, user, stream, space)
}

func (app *App) GetMessagesForStream(c *Context, currentUserId string, streamId string) ([]*model.Message, *model.AppErr) {
	/* if !app.CheckIfRoomMember(currentUserId, streamId) {
		return nil, &model.AppErr{
			Msg:        "Not Room Member",
			StatusCode: http.StatusForbidden,
		}
	} */

	messages, err := app.getMessages(streamId)
	if err != nil {
		log.Println(err)
		return nil, &model.AppErr{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if len(messages) == 0 {
		return nil, &model.AppErr{
			Msg:        "Empty messages",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return messages, nil
}

func (a *App) GetDefaultStream(currentUserId string) (string, *model.AppErr) {
	stream, err := a.getDefaultStream(currentUserId)
	if err != nil {
		log.Println("err")
		return "", &model.AppErr{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if stream == nil {
		log.Println("nil stream")
		return "", &model.AppErr{
			Msg:        "No default stream",
			StatusCode: http.StatusNotFound,
		}
	}

	return stream.ID, nil
}

func (a *App) createMessage(msg *model.Message) (*model.Message, error) {
	msg.PreSave()
	result := a.Store.DB.Create(&msg)
	return msg, result.Error
}

func (a *App) getMessages(streamId string) ([]*model.Message, error) {
	var messages []*model.Message
	result := a.Store.DB.Where("stream_id = ?", streamId).Find(&messages)
	return messages, result.Error
}

func (a *App) getDefaultStream(userId string) (*model.Stream, error) {
	var stream *model.Stream
	//result := a.Store.DB.Raw("SELECT * FROM messages WHERE author_id = ? ORDER BY created_at limit 1", userId).Scan(&stream)
	result := a.Store.DB.Table("streams").First(&stream)
	return stream, result.Error
}
