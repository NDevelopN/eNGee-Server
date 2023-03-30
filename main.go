package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	server "Engee-Server/server"
)

var mux = map[string]func(http.ResponseWriter, *http.Request){
	"/":  def,
	"/s": socket,
}

var upgrader = websocket.Upgrader{}

func main() {
	log.Println("Welcome to Engee-Server!")
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

func socket(writer http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Websocket connection established")
	defer conn.Close()

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		input := string(message)

		err = conn.WriteMessage(mt, []byte(input))
		if err != nil {
			log.Println("write failed:", err)
			break
		}

	}

}
