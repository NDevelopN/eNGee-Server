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

func loadEnv() {
	port, found := os.LookupEnv("SERVER_INNER")
	if found {
		config.Server.Port = port

		config.Database.Host = os.Getenv("POSTGRES_HOST")
		config.Database.Port = os.Getenv("POSTGRES_OUTER")
		config.Database.User = os.Getenv("POSTGRES_USER")
		config.Database.Pass = os.Getenv("POSTGRES_PASSWORD")
		config.Database.Name = os.Getenv("POSTGRES_DB")
	} else {
		log.Printf("Environment doesn't seem to be set. Checking config file.")
		loadConfig()
	}
}

func loadConfig() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Printf("Could not read config file: %v", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Could not read config data: %v", err)
	}

}

func main() {
	loadEnv()
	database.InitDB(config)
	handlers.Init()
	server.Serve(config)
}
