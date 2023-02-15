package app

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/RobleDev498/spaces/model"
	"gorm.io/gorm"
)

func (app *App) AddFriend(friend *model.Friend) (*model.User, *model.AppErr) {
	relation, err := app.addFriend(friend)
	if err != nil {
		return nil, &model.AppErr{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
			AppErrCode: model.ERROR_INTERNAL_ERROR,
		}
	}

	return app.GetUser(relation.FriendID)
}

func (app *App) GetFriends(userId string) ([]*model.User, *model.AppErr) {
	friends, err := app.getFriends(userId)
	if len(friends) == 0 {
		err = gorm.ErrEmptySlice
	}

	if err != nil {
		if errors.Is(err, gorm.ErrEmptySlice) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &model.AppErr{
				Msg:        "friends not found",
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

	return friends, nil
}

func (app *App) SanitizeFriends(friends []*model.User) {
	for _, friend := range friends {
		friend.SanitizeFriend()
	}
}

func (app *App) addFriend(friend *model.Friend) (*model.Friend, error) {
	friend.PreSave()
	result := app.Store.DB.Table("friends").Create(&friend)
	return friend, result.Error
}

func (app *App) getFriends(userId string) ([]*model.User, error) {
	var friends []*model.User
	result := app.Store.DB.Raw(`SELECT u.* FROM users u
	JOIN friends f ON f.user_id = u.id
	WHERE u.id = @userId UNION SELECT u.* FROM users u
	JOIN friends f ON f.friend_id = u.id
	WHERE u.id = @userId`, sql.Named("userId", userId)).Scan(&friends)

	return friends, result.Error
}

func (app *App) getFriend(userId, friendId string) (*model.Friend, error) {
	if userId > friendId {
		temp := userId
		userId = friendId
		friendId = temp
	}

	var relation *model.Friend
	result := app.Store.DB.Table("friends").Where("user_id = ? AND friend_id = ?", userId, friendId).First(&relation)

	return relation, result.Error
}
