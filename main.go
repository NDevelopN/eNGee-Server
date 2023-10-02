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
	port, found := os.LookupEnv("SPORT")
	if found {
		config.Server.Port = port

		config.Database.Host = os.Getenv("DHOST")
		config.Database.Host = os.Getenv("DPORT")
		config.Database.Host = os.Getenv("DUSER")
		config.Database.Host = os.Getenv("DPASS")
		config.Database.Host = os.Getenv("DNAME")
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
