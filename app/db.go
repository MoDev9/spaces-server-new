package app

import (
	"log"
	"os"
	"time"

	"github.com/RobleDev498/spaces/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (a *App) InitDb() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	dsn := a.Config.SqlSettings.DSN

	//dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Panic(err)
	}

	if err = sqlDB.Ping(); err != nil {
		log.Panic(err)
	}
	//return db
	a.SetDB(db)
}

func (app *App) Migrate() {
	db := app.Store.DB
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Message{})
	db.AutoMigrate(&model.Stream{})
	db.AutoMigrate(&model.Space{})
}
