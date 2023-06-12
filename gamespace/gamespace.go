package gamespace

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	db "Engee-Server/database"
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

	//Check if player exists
	plr, err := db.GetPlayer(msg.PID)
	if err != nil {
		log.Printf("[Error] Failed to get player from db in gsHandler: %v", err)
		u.SockSend(conn, "Error", msg.PID, msg.GID, "Failed to get player from db in gsHandler")
		return
	}

	//Check if game exists
	gm, err := db.GetGame(msg.GID)
	if err != nil {
		log.Printf("[Error] Failed to get game from db in gsHandler: %v", err)
		u.SockSend(conn, "Error", msg.PID, msg.GID, "Failed to get game from db in gsHandler")
		return
	}

	//Check if player is in the game
	if plr.GID != msg.GID {
		log.Printf("[Error] Player GID does not match game GID: %v", err)
		u.SockSend(conn, "Error", msg.PID, msg.GID, "Player GID does not match game GID")
		return
	}

	//Check if the player is the game leader
	leader := (msg.PID == gm.Leader)

	switch msg.Type {
	//Case for first connection
	case "Connect":
		//TODO this is not clear
		if Connect(conn, gm, plr) {
			gHandler(
				u.GameMsg{
					Type:    "Create",
					PID:     msg.PID,
					GID:     msg.GID,
					Content: gm.AdditionalRules,
				},
				Broadcast,
			)

		}
		gHandler(msg, Broadcast)
	case "Status":
		ChangePlayerStatus(conn, msg.GID, plr, msg.Content, leader)

		//If autostart enabled, start the game after more than half of players are ready
		if true {
			ready := db.GetGamePReady(msg.GID)
			threshold := db.GetGamePCount(msg.GID) / 2
			if ready > threshold {
				Start(conn, gm)
				msg.Type = "Start"
				gHandler(msg, Broadcast)
			}
		}
	case "Leave":
		gHandler(msg, Broadcast)
		Leave(conn, plr, gm)
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
		gHandler(msg, Broadcast)
	case "End":
		if !leader {
			u.SockSend(conn, "Error", msg.GID, msg.PID, "Player is not the leader")
			return
		}
		End(conn, gm)
		gHandler(msg, Broadcast)
	case "Restart":
		if !leader {
			u.SockSend(conn, "Error", msg.GID, msg.PID, "Player is not the leader")
			return
		}
		Restart(conn, gm)
		gHandler(msg, Broadcast)
	case "Remove":
		if !leader {
			u.SockSend(conn, "Error", msg.GID, msg.PID, "Player is not the leader")
			return
		}
		Remove(conn, gm, msg.Content)
	case "Rules":
		if !leader {
			u.SockSend(conn, "Error", msg.GID, msg.PID, "Player is not the leader")
			return
		}
		UpdateRules(conn, gm, msg.Content)
		gHandler(msg, Broadcast)
	default:
		gHandler(msg, Broadcast)
	}
}

func UpdateStatus(conn *websocket.Conn, gm u.Game) {
	gMsg := u.GameMsg{
		Type:    "Status",
		GID:     gm.GID,
		PID:     "",
		Content: gm.Status,
	}

	msg, err := json.Marshal(gMsg)
	if err != nil {
		u.SockSend(conn, "Error", "", "", "Failed to send status update message")
		log.Printf("[Error] Failed to marshal status update message: %v", err)
		return
	}

	Broadcast(gm.GID, msg)
}

func UpdatePlayerList(gid string) {
	plrs, err := db.GetGamePlayers(gid)
	if err != nil {
		log.Printf("[Error] Failed to get player list from game %v", err)
		return
	}
	list, err := json.Marshal(plrs)
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

	Broadcast(gid, msg)
}

func GameSpace(w http.ResponseWriter, r *http.Request, gHandler u.GHandler) {
	u.Sock(w, r, handler, gHandler)
}

func Broadcast(gid string, msg []byte) {
	plrs, err := db.GetGamePlayers(gid)
	if err != nil {
		log.Printf("[Error] Failed to get players for broadcast: %v", err)
		return
	}

	for _, p := range plrs {
		u.Connections[gid][p.PID].WriteMessage(websocket.TextMessage, msg)
	}

}
