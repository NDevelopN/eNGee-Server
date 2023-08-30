package main

import (
	"Engee-Server/database"
	"Engee-Server/handlers"
	"Engee-Server/server"
)

func main() {
	database.InitDB()
	handlers.Init()
	server.Serve()
}
