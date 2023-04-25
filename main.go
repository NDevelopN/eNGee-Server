package main

import (
	"log"
	"time"

	server "Engee-Server/server"
)

func main() {
	log.Println("Welcome to Engee-Server!")
	go func() {
		server.Serve()
	}()

	for {
		time.Sleep(1 * time.Second)
	}
}
