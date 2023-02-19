package main

import (
	"github.com/labstack/gommon/log"
	"my-rest-api/config"
	"my-rest-api/internal/app"
	"my-rest-api/internal/service"
	"my-rest-api/pkg/store"
)

func main() {

	cfg := config.NewConfig()

	log.Info("Configuration completed!")

	storage, err := store.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer func(storage *store.Store) {
		err := storage.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(storage)

	log.Info("Connection to database is succeeded")

	userService := service.NewUserService(storage)
	log.Info("Service created successfully")

	api := app.NewApi(cfg, userService)
	log.Fatal(api.Start())

}
