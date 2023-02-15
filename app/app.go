package app

import (
	"github.com/RobleDev498/spaces/config"
	"github.com/RobleDev498/spaces/store"
	"gorm.io/gorm"
)

type App struct {
	db    *gorm.DB
	srv   *Server
	Store *store.Store

	Config *config.Config
}

func New(s *Server, config *config.Config) *App {
	a := &App{
		//db:  db,
		Config: config,
		srv:    s,
		Store:  &store.Store{},
	}

	//a.srv.hub = a.NewHub()
	return a
}

func (a *App) Init() {
	a.InitDb()
}

func (a *App) Srv() *Server {
	return a.srv
}

func (a *App) SetDB(db *gorm.DB) {
	a.Store.DB = db
}
