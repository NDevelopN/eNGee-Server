package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	server "Engee-Server/server"
)

var mux = map[string]func(http.ResponseWriter, *http.Request){
	"/": def,
}

func main() {
	go func() {
		server.Serve(mux)
	}()

	for {
		time.Sleep(1 * time.Second)
	}
}

// default function, return error, but log request
func def(writer http.ResponseWriter, request *http.Request) {
	reqBody, _ := ioutil.ReadAll(request.Body)
	log.Print("Request to \"/\": \n")
	log.Print(reqBody)
	log.Print("\n")

	http.Error(writer, "Invalid request", http.StatusBadRequest)
}
