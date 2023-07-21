package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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

	game, err := g.GetGame(user.GID)
	if err != nil {
		http.Error(w, "Failed to get matching game", http.StatusBadRequest)
		log.Printf("[Error] getting game with requested GID: %v", err)
		return
	}

	conn, err := upgradeConnection(w, r)
	if err != nil {
		http.Error(w, "Failed to upgrade to websocket connection", http.StatusInternalServerError)
		log.Printf("[Error] upgrading connection: %v", err)
		return
	}

	conn.SetCloseHandler(handleClose)

	err = utils.AddConnection(user.GID, user.UID, conn)
	if err != nil {
		http.Error(w, "Failed to add connection to pool", http.StatusInternalServerError)
		log.Printf("[Error] adding connection to pool: %v", err)
	}

	gInfo, err := json.Marshal(game)
	if err != nil {
		http.Error(w, "Failed to marshal game information", http.StatusInternalServerError)
		log.Printf("[Error] Failed to game info: %v", err)
		return
	}

	msg := utils.GameMsg{
		Type:    "Info",
		UID:     user.UID,
		GID:     user.GID,
		Content: string(gInfo),
	}

	reply, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, "Failed to marshal welcome message", http.StatusInternalServerError)
		log.Printf("[Error] Failed to marshal message: %v", err)
		return
	}

	conn.WriteMessage(websocket.TextMessage, reply)

	go handleIncoming(user.GID, user.UID)

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

	return conn, nil
}

func handleClose(code int, text string) error {
	if code == websocket.CloseNoStatusReceived {
		return fmt.Errorf("connnection closed: without status")
	}

	return fmt.Errorf("connection closed: %v", text)
}

func handleIncoming(gid string, uid string) {
	pool, err := utils.GetConnections(gid)
	if err != nil {
		log.Printf("[Error] getting connection for handler: %v", err)
	}
	conn := pool[uid]

	for utils.CheckConnection(gid, uid) {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			if strings.Split(err.Error(), ":")[0] == "connection closed" {
				log.Printf("[Close] connection closed: %v", err)
				utils.RemoveConnection(gid, uid)
				msg := utils.GameMsg{
					Type: "Leave",
					GID:  gid,
					UID:  uid,
				}

				_, err = gs.GamespaceHandle(msg)
				if err != nil {
					conn.Close()
				}
				return
			} else {
				log.Printf("[Error] Failed to read message: %v", err)
				continue
			}
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
