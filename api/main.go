package main

import (
	"api/cache"
	"api/cmd"
	"api/config"
	"api/database"
	apilogger "api/logger"
	"api/models"
	"api/server"
)

func main() {

	// Setup logger
	apilogger.SetLogLevel(config.Log().Level)

	// Setup routes
	db, err := database.GetDB()
	if err != nil {
		panic(err)
	}

	cache, err := cache.NewReddis()
	if err != nil {
		panic(err)
	}

	defer db.Close()
	defer cache.Close()

	// initial fill database with initial data with fill-db command if the arguments are more than 1
	cmd.NewCmd().Execute()

	configServer := models.NewServerConfig(
		db.DB,
		config.Server().Port,
		cache,
	)

	server := server.NewServer(configServer)
	err = server.Start()
	if err != nil {
		apilogger.Logger().Err(err).Msg("Server error")
		panic(err)
	}

}
