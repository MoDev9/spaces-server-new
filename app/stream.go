package app

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/RobleDev498/spaces/model"
	"gorm.io/gorm"
)

func (app *App) CreateStream(currentUserId string, stream *model.Stream) (*model.Stream, *model.AppErr) {
	if stream.SpaceID == "" {
		stream.SpaceID = "h"
	}

	stream, err := app.Store.CreateStream(stream)
	if err != nil {
		return nil, &model.AppErr{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
			AppErrCode: model.ERROR_INTERNAL_ERROR,
		}
	}

	if (stream.Type == model.DM || stream.Type == model.GROUP_DM) && stream.OwnerID == currentUserId {
		_, err := app.AddUserToRoom(currentUserId, &model.RoomMember{
			StreamID: stream.ID,
			UserID:   currentUserId,
		}, nil)

		if err != nil {
			return nil, err
		}
	}

	return stream, nil
}

func (app *App) CheckIfRoomMember(userId string, roomId string) bool {
	rooms, err := app.Store.GetRooms(userId)

	if len(rooms) == 0 {
		return false
	}

	if err != nil {
		return false
	}

	for _, room := range rooms {
		if room.ID == roomId {
			return true
		}
	}

	return false
}

func (app *App) AddUserToRoom(currentUserId string, roomMember *model.RoomMember, user *model.User) (*model.RoomMember, *model.AppErr) {
	stream, appErr := app.GetStream(roomMember.StreamID)
	if appErr != nil {
		return nil, appErr
	}

	if stream.OwnerID != currentUserId {
		return nil, &model.AppErr{
			Msg:        "User does not have permission to add friend to room",
			StatusCode: http.StatusInternalServerError,
			AppErrCode: model.ERROR_INTERNAL_ERROR,
		}
	}

	rm, err := app.Store.AddFriendToRoom(roomMember)
	if err != nil {
		return nil, &model.AppErr{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
			AppErrCode: model.ERROR_INTERNAL_ERROR,
		}
	}

	if user != nil {
		if stream.Name == "" {
			stream.Name = user.Username
		} else {
			stream.Name = fmt.Sprintf("%s, %s", stream.Name, user.Username)
		}
	}

	_, err = app.Store.UpdateStream(stream)
	if err != nil {
		return nil, &model.AppErr{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
			AppErrCode: model.ERROR_INTERNAL_ERROR,
		}
	}

	return rm, nil
}

func (app *App) AddFriendsToRoom(roomMembers []*model.RoomMember) ([]*model.RoomMember, *model.AppErr) {
	rm, err := app.Store.AddFriendsToRoom(roomMembers)
	if err != nil {
		return nil, &model.AppErr{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
			AppErrCode: model.ERROR_INTERNAL_ERROR,
		}
	}

	return rm, nil
}

func (app *App) GetStream(id string) (*model.Stream, *model.AppErr) {
	stream, err := app.Store.GetStream(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &model.AppErr{
				Msg:        "stream not found",
				StatusCode: http.StatusNotFound,
				AppErrCode: model.ERROR_RECORD_NOT_FOUND,
			}
		}

		return nil, &model.AppErr{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
			AppErrCode: model.ERROR_INTERNAL_ERROR,
		}
	}
	stream.Sanitize()
	return stream, nil
}

func (app *App) GetStreamMembers(stream_id string) ([]*model.User, *model.AppErr) {
	members, err := app.Store.GetStreamMembers(stream_id)
	if len(members) == 0 {
		err = gorm.ErrRecordNotFound
	}

	if err != nil {
		if errors.Is(err, gorm.ErrEmptySlice) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &model.AppErr{
				Msg:        "stream members not found",
				StatusCode: http.StatusNotFound,
				AppErrCode: model.ERROR_RECORD_NOT_FOUND,
			}
		}

		return nil, &model.AppErr{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
			AppErrCode: model.ERROR_INTERNAL_ERROR,
		}
	}

	return members, nil
}

func (app *App) GetRooms(userId string) ([]*model.Stream, *model.AppErr) {
	streams, err := app.Store.GetRooms(userId)

	if len(streams) == 0 {
		err = gorm.ErrEmptySlice
	}

	if err != nil {
		if errors.Is(err, gorm.ErrEmptySlice) || errors.Is(err, gorm.ErrRecordNotFound) {
			/* return nil, &model.AppErr{
				Msg:        "streams not found",
				StatusCode: http.StatusNotFound,
				AppErrCode: model.ERROR_RECORD_NOT_FOUND,
			} */
			emptySlice := make([]*model.Stream, 0)
			return emptySlice, nil
		}

		return nil, &model.AppErr{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
			AppErrCode: model.ERROR_INTERNAL_ERROR,
		}
	}

	for _, stream := range streams {
		streamMembers, _ := app.Store.GetStreamMembers(stream.ID)
		stream.Members = streamMembers
	}

	return streams, nil
}

func (app *App) SanitizeRooms(session model.Session, streams []*model.Stream) {
	//For now
	for _, stream := range streams {
		stream.Sanitize()
	}
}
