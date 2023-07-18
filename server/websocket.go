package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	g "Engee-Server/game"
	u "Engee-Server/user"
	utils "Engee-Server/utils"

	gs "Engee-Server/gamespace"
)

func Connect(c *gin.Context) {
	w := c.Writer
	r := c.Request

	uid, err := GetID(c)
	if err != nil {
		http.Error(w, "Could not get user ID from request path", http.StatusBadRequest)
		log.Printf("[Error] Getting user ID: %v", err)
		return
	}

	user, err := u.GetUser(uid)
	if err != nil {
		http.Error(w, "Failed to get user with matching ID", http.StatusBadRequest)
		log.Printf("[Error] Getting User: %v", err)
		return
	}

	// Check if the provided GID matches an existing game

	_, err = g.GetGame(user.GID)
	if err != nil {
		http.Error(w, "Failed to get matching game", http.StatusBadRequest)
		log.Printf("[Error] getting game with requested GID: %v", err)
		return
	}

	conn, err := upgradeConnection(w, r)
	if err != nil {
		log.Printf("[Error] upgrading connection: %v", err)
		return
	}

	err = utils.AddConnection(user.GID, user.UID, conn)
	if err != nil {
		log.Printf("[Error] adding connection to pool: %v", err)
	}
}

func upgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {

	var upgrader = websocket.Upgrader{}

	// Create websocket connection
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade websocket connection", http.StatusInternalServerError)
		return nil, fmt.Errorf("failed to create a websocket connection: %v", err)
	}

	// Maintain the connection
	go handleIncoming(conn)

	return conn, nil
}

func handleIncoming(conn *websocket.Conn) {
	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[Error] Failed to read message: %v", err)
			continue
		}

		if messageType != websocket.TextMessage {
			log.Printf("[Error] Received unexpected message type: %v", messageType)
			continue
		}

		var msg = utils.GameMsg{}

		err = json.Unmarshal(data, &msg)
		if err != nil {
			log.Printf("[Error] Cannot unmarshal received message: %v", err)
			continue
		}

		log.Printf("[Request] Received message: %v", msg)

		msg, err = gs.GamespaceHandle(msg)
		if err != nil {
			log.Printf("[Error] Handling message: %v", err)
		}

		reply, err := json.Marshal(msg)
		if err != nil {
			log.Printf("[Error] Failed to marshal message: %v", err)
			continue
		}

		conn.WriteMessage(websocket.TextMessage, reply)
	}
}
