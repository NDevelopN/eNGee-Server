package main

import (
	"Engee-Server/database"
	"Engee-Server/handlers"
	"Engee-Server/server"
	"Engee-Server/utils"
)

func main() {
	config := utils.LoadEnv()
	database.InitDB(config)
	handlers.Init()
	server.Serve(config)
}
