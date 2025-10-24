package main

import (
	"api/cache"
	"api/cmd"
	"api/config"
	"api/database"
	apilogger "api/logger"
	"api/models"
	"api/server"
	"fmt"
)

func main() {

	// Setup logger
	apilogger.InitLogger()
	apilogger.SetLogLevel(config.Log().Level)

	// Setup routes
	db, err := database.GetDB()
	if err != nil {
		panic(fmt.Errorf("error getting db: %v", err))
	}

	cache, err := cache.NewReddis()
	if err != nil {
		panic(fmt.Errorf("error getting cache db: %v", err))
	}

	defer db.Close()
	defer cache.Close()

	// initial fill database with initial data with fill-db command if the arguments are more than 1
	err = cmd.NewCmd().Execute()
	if err != nil {
		apilogger.Logger().Err(err).Msg("Error executing command")
		panic(err)
	}

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
