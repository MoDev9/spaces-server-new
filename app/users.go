package app

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/RobleDev498/spaces/model"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

func (app *App) CreateUserFromRegister(user *model.User) (*model.User, *model.AppErr) {
	appErr := app.ValidUsername(user.Username)
	if appErr != nil {
		return nil, appErr
	}

	appErr = IsPasswordValid(user.Password)
	if appErr != nil {
		return nil, appErr
	}

	ruser, err := app.Store.CreateUser(user)
	if err != nil {
		if pgErr := err.(*pgconn.PgError); pgErr != nil {
			if pgErr.Code == model.ERROR_PG_UNIQUE {
				constraint_name := pgErr.ConstraintName

				switch constraint_name {
				case model.ERROR_CONSTRAINT_USERNAME:
					appErr = &model.AppErr{
						Msg:        fmt.Sprintf("Username %s taken", user.Username),
						AppErrCode: model.ERROR_USERNAME_TAKEN,
						StatusCode: http.StatusConflict,
					}
				case model.ERROR_CONSTRAINT_EMAIL:
					appErr = &model.AppErr{
						Msg:        fmt.Sprintf("User with email %s already exists", user.Email),
						AppErrCode: model.ERROR_USERNAME_TAKEN,
						StatusCode: http.StatusConflict,
					}
				default:
					appErr = nil
				}
				return nil, appErr
			}
		}
	}

	ruser.Sanitize()
	return ruser, nil
}

func (app *App) GetUserByUsername(username string) (*model.User, *model.AppErr) {
	user, err := app.Store.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &model.AppErr{
				Msg:        "user not found",
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
	user.Sanitize()
	return user, nil
}

func (app *App) GetUser(id string) (*model.User, *model.AppErr) {
	user, err := app.Store.GetUser(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &model.AppErr{
				Msg:        "user not found",
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
	user.Sanitize()
	return user, nil
}

func (a *App) UserCanBeViewed(userId, otherUserId string) (bool, *model.AppErr) {
	if userId == otherUserId {
		return true, nil
	}

	return false, nil
}

/* func (app *App) deleteUser(id uint) error {
	return app.Store.DB.Delete(u, id).Error
} */
