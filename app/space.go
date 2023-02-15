package app

import (
	"errors"
	"log"
	"net/http"

	"github.com/RobleDev498/spaces/model"
	"gorm.io/gorm"
)

func (app *App) CreateSpace(owner, name string, permissions map[string]interface{}) {

}

func (app *App) GetSpacesForUser(userId string) ([]*model.Space, *model.AppErr) {

	spaces, err := app.getSpaces(userId)
	if err != nil {
		log.Println(err)
		if errors.Is(err, gorm.ErrEmptySlice) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &model.AppErr{
				Msg:        "spaces not found. Error: " + err.Error(),
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

	return spaces, nil
}

func (app *App) GetSpace(spaceId string) (*model.Space, *model.AppErr) {
	/* u64, err := strconv.ParseUint(spaceId, 10, 32)
	if err != nil {
		fmt.Println(err)
		return nil, &model.AppErr{
			Msg:        "space_id must be of valid format",
			StatusCode: http.StatusUnprocessableEntity,
		}
	} */

	space, err := app.getSpace(spaceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &model.AppErr{
				Msg:        "space not found",
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

	return space, nil
}

func (app *App) getSpace(spaceId string) (*model.Space, error) {
	var space *model.Space
	result := app.Store.DB.First(&space, spaceId)
	return space, result.Error
}

func (app *App) getSpaces(userId string) ([]*model.Space, error) {
	var spaces []*model.Space
	result := app.Store.DB.Raw(`SELECT s.* FROM spaces s 
		JOIN space_users su ON s.id == su.space_id 
		JOIN users u ON su.user_id = u.id 
		WHERE u.id = ?`,
		userId).Scan(&spaces)
	return spaces, result.Error
}
