package main

import (
	"Engee-Server/config"
	"Engee-Server/server"
)

func main() {
	config := config.ReadConfig()
	server.Serve(config.Server.Port)
}
