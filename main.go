package main

import (
	"Engee-Server/database"
	"Engee-Server/handlers"
	"Engee-Server/server"
	"Engee-Server/utils"
	"encoding/json"
	"log"
	"os"
)

var config utils.Config

func loadConfig() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Could not read config file: %v", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Could not read config data: %v", err)
	}
}

func main() {
	loadConfig()
	database.InitDB(config)
	handlers.Init()
	server.Serve(config)
}
