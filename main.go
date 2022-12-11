package main

import (
	"fmt"
	"log"

	"github.com/letenk/pokedex/router"
	"github.com/letenk/pokedex/util"
)

func main() {
	// Load Config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	fmt.Println("Server is starting...")
	// Open connection to postgres
	db := util.SetupDB(config.DB_SOURCE)

	// Setup Router
	router := router.SetupRouter(db)
	app_port := fmt.Sprintf(":%s", config.APP_PORT)
	router.Run(app_port)
}
