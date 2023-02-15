package app

import (
	"errors"
	"net/http"

	"github.com/RobleDev498/spaces/model"
	"gorm.io/gorm"
)

func (app *App) RevokeSessionById(sessionId string) *model.AppErr {
	session, err := app.GetSessionById(sessionId)
	if err != nil {
		return err
	}

	return app.RevokeSession(session)
}

func (app *App) RevokeSession(session *model.Session) *model.AppErr {
	err := app.Store.RemoveSession(session)
	if err != nil {
		return &model.AppErr{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
			AppErrCode: model.ERROR_INTERNAL_ERROR,
		}
	}

	return nil
}

func (app *App) GetSession(token string) (*model.Session, *model.AppErr) {
	session, err := app.Store.GetSession(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &model.AppErr{
				Msg:        "session not found",
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

	return session, nil
}

func (app *App) GetSessionById(id string) (*model.Session, *model.AppErr) {
	session, err := app.Store.GetSessionById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &model.AppErr{
				Msg:        "session not found",
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

	return session, nil
}
