package gamespace

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"

	u "Engee-Server/utils"
)

// TODO freeze timer?
func Pause(conn *websocket.Conn, gm u.Game) {
	if gm.Status == "Pause" {
		gm.Status = gm.OldStatus
		UpdateStatus(conn, gm.GID, gm.Status)
	} else {
		gm.OldStatus = gm.Status
		gm.Status = "Pause"
		UpdateStatus(conn, gm.GID, gm.Status)
	}

	u.Games[gm.GID] = gm
}

// TODO add game spec
func Start(conn *websocket.Conn, gm u.Game) {
	//TODO add option to toggle this as requirement
	if len(gm.Players) >= gm.Rules.MinPlrs {
		UpdateStatus(conn, gm.GID, "Play")
		//TODO: here goes the start function
	} else {
		u.SockSend(conn, "Block", "", "", "There are not enough players in the game.")
		log.Printf("[Block] Attempted to start without enough players in the game.")
	}

	u.Games[gm.GID] = gm
}

func End(conn *websocket.Conn, gm u.Game) {
	//TODO
}

func Restart(conn *websocket.Conn, gm u.Game) {
	UpdateStatus(conn, gm.GID, "Restart")

	for i := range gm.Players {
		gm.Players[i].Status = "New"
	}

	//TODO: Remove
	time.Sleep(1 * time.Second)
	UpdateStatus(conn, gm.GID, "Lobby")

	u.Games[gm.GID] = gm
}

func UpdateRules(conn *websocket.Conn, gm u.Game, content string) {
	var rules u.Game
	err := json.Unmarshal([]byte(content), &rules)
	if err != nil {
		u.SockSend(conn, "Error", "", "", "Could not read rules update message")
		log.Printf("[Error] Failed to unmarshal new rules: %v", err)
		return
	}

	change := false
	if gm.Type != rules.Type {
		gm.Type = rules.Type
		change = true
	}

	if gm.Rules != rules.Rules {
		gm.Rules = rules.Rules
		change = true
	}

	if !change {
		log.Printf("No changes to game rules provided")
		return
	}

	//Put game in restarting state
	Restart(conn, gm)

	gUpdate, err := json.Marshal(gm)
	if err != nil {
		log.Printf("[Error] Failed to marshal game rules for update: %v", err)
		return
	}

	gMsg := u.GameMsg{
		Type:    "Update",
		PID:     "",
		GID:     gm.GID,
		Content: string(gUpdate),
	}

	msg, err := json.Marshal(gMsg)
	if err != nil {
		log.Printf("[Error] Failed to marshal update message: %v", err)
		return
	}

	// Send the players the rules update
	u.Broadcast(gm.GID, msg)
	time.Sleep(1 * time.Second)

	// Leave restart state, return to lobby
	UpdateStatus(conn, gm.GID, "Lobby")
	u.Games[gm.GID] = gm
}

func Remove(conn *websocket.Conn, gm u.Game, content string) {
	//TODO
}
