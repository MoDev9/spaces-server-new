package main

import (
	"log"

	"github.com/RobleDev498/spaces/api"
	"github.com/RobleDev498/spaces/app"
	"github.com/RobleDev498/spaces/config"
)

func main() {
	configPath := "/etc/whiteboard/config.yml"
	err := config.ValidateConfigPath(configPath)
	if err != nil {
		log.Fatal(err)
	}

	config, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	s, err := app.NewServer()
	s.Config = config
	if err != nil {
		log.Panic(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	a := app.New(s, config)
	//a.StartHub()
	a.Init()
	api.Init(a, s.RootRouter)

	hub := a.NewHub()
	go hub.Start()
	a.SetHub(hub)
	err = s.Start()
	if err != nil {
		log.Panic(err)
	}
}
