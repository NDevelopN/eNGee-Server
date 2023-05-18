package gamespace

import (
	u "Engee-Server/utils"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func Connect(conn *websocket.Conn, msg u.GameMsg) bool {
	first := false

	// Add new connection to map
	gCon := u.Connections[msg.GID]
	gCon[msg.PID] = conn
	u.Connections[msg.GID] = gCon

	gm := u.Games[msg.GID]
	if gm.Leader == "" || gm.Leader == msg.PID {
		gm.Leader = msg.PID
		first = true
	}

	info, err := json.Marshal(gm)
	if err != nil {
		log.Printf("[Error] Failed to marshal game info: %v", err)
		return false
	}

	u.SockSend(conn, "Info", msg.GID, msg.PID, string(info))

	UpdatePlayerList(msg.GID)
	return first
}

func Status(conn *websocket.Conn, gm u.Game, pid string, status string) bool {
	ready := 0

	for i, p := range gm.Players {
		if p.PID == pid {
			p.Status = status
			gm.Players[i] = p
			UpdatePlayerList(gm.GID)
		}

		if p.Status == "Ready" {
			ready++
		}
	}

	return ready > (len(gm.Players) / 2)
}

func Leave(conn *websocket.Conn, msg u.GameMsg) {
	leader := u.RemovePlayer(msg.GID, msg.PID)
	if leader != "" {
		u.SockSend(u.Connections[msg.GID][leader], "Leader", msg.GID, leader, "")
		UpdatePlayerList(msg.GID)
	}
}
