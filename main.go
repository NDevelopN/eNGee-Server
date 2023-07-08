package main

import (
	database "Engee-Server/database"
	server "Engee-Server/server"
)

func main() {
	database.InitDB()
	server.Serve()
}
