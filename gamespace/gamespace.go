package gamespace

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	u "Engee-Server/utils"
)

var handler u.MHandler = func(conn *websocket.Conn, data []byte, gHandler u.GHandler) {

	var msg u.GameMsg
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("[Error] Failed to unmarshal message: %v", err)
		u.SockSend(conn, "Error", "", "", "Failed to read message")
		return
	}

	found, _ := u.CheckForPlayer(msg.PID)
	if !found {
		u.SockSend(conn, "Error", "", "", "Player was not found")
		return
	}

	found, gm := u.CheckForGame(msg.GID)
	if !found {
		u.SockSend(conn, "Error", "", "", "Game was not found")
		return
	}

	found = u.CheckGameForPlayer(gm, msg.PID)
	if !found {
		u.SockSend(conn, "Error", "", "", "Player was not found in the game")
		return
	}

	leader := (msg.PID == gm.Leader)

	switch msg.Type {
	case "Connect":
		if Connect(conn, msg) {
			gHandler(
				u.GameMsg{
					Type:    "Create",
					PID:     msg.PID,
					GID:     msg.GID,
					Content: gm.Rules.Additional,
				},
			)
		}
		gHandler(msg)
	case "Status":
		Status(conn, gm, msg.PID, msg.Content)
	case "Leave":
		gHandler(msg)
		Leave(conn, msg)
	case "Pause":
		if !leader {
			u.SockSend(conn, "Error", msg.GID, msg.PID, "Player is not the leader")
			return
		}
		Pause(conn, gm)
	case "Start":
		if !leader {
			u.SockSend(conn, "Error", msg.GID, msg.PID, "Player is not the leader")
			return
		}
		Start(conn, gm)
		gHandler(msg)
	case "End":
		if !leader {
			u.SockSend(conn, "Error", msg.GID, msg.PID, "Player is not the leader")
			return
		}
		End(conn, gm)
	case "Restart":
		if !leader {
			u.SockSend(conn, "Error", msg.GID, msg.PID, "Player is not the leader")
			return
		}
		Restart(conn, gm)
	case "Remove":
		if !leader {
			u.SockSend(conn, "Error", msg.GID, msg.PID, "Player is not the leader")
			return
		}
		Remove(conn, gm, msg.Content)
	default:
		gHandler(msg)
		//TODO call game handler
	}
}

func UpdateStatus(conn *websocket.Conn, gid string, status string) {
	gMsg := u.GameMsg{
		Type:    "Status",
		GID:     gid,
		PID:     "",
		Content: status,
	}

	msg, err := json.Marshal(gMsg)
	if err != nil {
		u.SockSend(conn, "Error", "", "", "Failed to send status update message")
		log.Printf("[Error] Failed to marshal status update message: %v", err)
		return
	}

	u.Broadcast(gid, msg)
}

func UpdatePlayerList(gid string) {
	list, err := json.Marshal(u.PlrList{Players: u.Games[gid].Players})
	if err != nil {
		log.Printf("[Error] Failed to marshal player list: %v", err)
		return
	}

	gMsg := u.GameMsg{
		Type:    "Players",
		GID:     gid,
		PID:     "",
		Content: string(list),
	}

	msg, err := json.Marshal(gMsg)
	if err != nil {
		log.Printf("[Error] Failed to marshal player list update message: %v", err)
		return
	}

	u.Broadcast(gid, msg)
}

func GameSpace(w http.ResponseWriter, r *http.Request, gHandler u.GHandler) {
	u.Sock(w, r, handler, gHandler)
}
