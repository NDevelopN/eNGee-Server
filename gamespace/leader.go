package gamespace

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"

	db "Engee-Server/database"
	u "Engee-Server/utils"
)

// TODO freeze timer?
func Pause(conn *websocket.Conn, gm u.Game) {
	if gm.Status == "Pause" {
		gm.Status = gm.OldStatus
	} else {
		gm.OldStatus = gm.Status
		gm.Status = "Pause"
	}

	err := db.UpdateGame(gm)
	if err != nil {
		log.Printf("[Error] Could not update game pause status in database: %v", err)
		u.SockSend(conn, "Error", "", "", "Could not update game pause status in database")
		return
	}

	UpdateStatus(conn, gm)
}

// TODO add game spec
func Start(conn *websocket.Conn, gm u.Game) {
	//TODO add option to toggle this as requirement
	if true {
		plrCount := db.GetGamePCount(gm.GID)

		if plrCount < gm.MinPlrs {
			u.SockSend(conn, "Block", "", "", "There are not enough players in the game.")
			log.Printf("[Block] Attempted to start without enough players in the game.")
			return
		}
	}

	gm.Status = "Play"

	err := db.UpdateGame(gm)
	if err != nil {
		log.Printf("[Error] Could not update game status in database: %v", err)
		u.SockSend(conn, "Error", "", "", "Could not update game status in database")
		return
	}

	err = db.UpdateGamePlayerStatus(gm.GID, "In Game")
	if err != nil {
		log.Printf("[Error] Could not update player statuses in database: %v", err)
		u.SockSend(conn, "Error", "", "", "Could not update player status in database")
		return
	}

	UpdatePlayerList(gm.GID)

	UpdateStatus(conn, gm)
}

func End(conn *websocket.Conn, gm u.Game) {
	err := db.RemoveGame(gm.GID)
	if err != nil {
		log.Printf("[Error] Failed to remove game after last player left: %v", err)
	}

	gMsg := u.GameMsg{
		Type:    "End",
		GID:     gm.GID,
		PID:     "",
		Content: "",
	}

	msg, err := json.Marshal(gMsg)
	if err != nil {
		log.Printf("[Error] Failed to marshal game end message: %v", err)
		return
	}

	Broadcast(gm.GID, msg)
	delete(u.Connections, gm.GID)
}

func Restart(conn *websocket.Conn, gm u.Game) {
	gm.Status = "Restart"

	err := db.UpdateGame(gm)
	if err != nil {
		log.Printf("[Error] Could not update game in database: %v", err)
		u.SockSend(conn, "Error", "", "", "Could not update game in database")
		return
	}

	UpdateStatus(conn, gm)

	err = db.UpdateGamePlayerStatus(gm.GID, "Joined")
	if err != nil {
		log.Printf("[Error] Could not update players of game in database: %v", err)
		u.SockSend(conn, "Error", "", "", "Could not update players of game in database")
		return
	}

	gm.OldStatus = ""
	gm.Status = "Lobby"

	err = db.UpdateGame(gm)
	if err != nil {
		log.Printf("[Error] Could not update game in database: %v", err)
		u.SockSend(conn, "Error", "", "", "Could not update game in database")
		return
	}

	UpdatePlayerList(gm.GID)
	UpdateStatus(conn, gm)
}

func UpdateRules(conn *websocket.Conn, gm u.Game, content string) {

	var rules u.Game
	err := json.Unmarshal([]byte(content), &rules)
	if err != nil {
		log.Printf("[Error] Could not parse new rules: %v", err)
		u.SockSend(conn, "Error", "", "", "Could not parse new rules")
		return
	}

	Restart(conn, gm)

	//TODO some checks here

	gm.Type = rules.Type
	gm.MinPlrs = rules.MinPlrs
	gm.MaxPlrs = rules.MaxPlrs
	gm.AdditionalRules = rules.AdditionalRules

	gUpdate, err := json.Marshal(gm)
	if err != nil {
		log.Printf("[Error] Could not marshal rules update: %v", err)
		u.SockSend(conn, "Error", "", "", "Could not marshal updated rules")
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
		u.SockSend(conn, "Error", "", "", "Could not marshal update message")
		return
	}

	err = db.UpdateGame(gm)
	if err != nil {
		log.Printf("[Error] Could not update game in database: %v", err)
		u.SockSend(conn, "Error", "", "", "Could not update game in database")
		return
	}

	// Send the players the rules update
	Broadcast(gm.GID, msg)

}

func Remove(conn *websocket.Conn, gm u.Game, content string) {
	//TODO
}
