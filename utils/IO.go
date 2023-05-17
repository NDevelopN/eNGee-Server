package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func Extract[O *Player | *Game | *Join](r *http.Request, obj O) error {
	reqBody, err := io.ReadAll(r.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(reqBody, obj)
	return err
}

func PackSend[O Player | Game | GameInfo | ACK](w http.ResponseWriter, msg O, e string) error {
	response, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, e, http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, e, http.StatusInternalServerError)
	}

	return err
}

type MHandler func(*websocket.Conn, []byte)

func Sock(w http.ResponseWriter, r *http.Request, handleMessage MHandler) {
	var upgrader = websocket.Upgrader{}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade websocket connection", http.StatusInternalServerError)
		log.Printf("Could not create websocket connection: %v", err)
	}

	go func() {
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Failure in reading message from websocket: %v", err)
				return
			}

			if messageType != websocket.TextMessage {
				log.Printf("Got unexpected message type: %v", messageType)
				log.Printf("Wanted %v", websocket.TextMessage)
			}

			handleMessage(conn, p)
		}
	}()
}
