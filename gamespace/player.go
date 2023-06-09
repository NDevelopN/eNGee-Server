package gamespace

import (
	db "Engee-Server/database"
	u "Engee-Server/utils"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func Connect(conn *websocket.Conn, gm u.Game, plr u.Player) bool {
	first := false

	// Add new connection to map
	gCon := u.Connections[gm.GID]
	gCon[plr.PID] = conn
	u.Connections[gm.GID] = gCon

	plr.Status = "Joined"

	if gm.Leader == "" || gm.Leader == plr.PID {
		gm.Leader = plr.PID
		first = true
	}

	err := db.UpdatePlayer(plr)
	if err != nil {
		log.Printf("[Error] Failed to update player status: %v", err)
		u.SockSend(conn, "Error", gm.GID, plr.PID, "Failed to update player status")
		return false
	}

	info, err := json.Marshal(gm)
	if err != nil {
		log.Printf("[Error] Failed to marshal game info: %v", err)
		return false
	}

	u.SockSend(conn, "Info", gm.GID, plr.PID, string(info))

	UpdatePlayerList(gm.GID)
	return first
}

func ChangePlayerStatus(conn *websocket.Conn, gid string, plr u.Player, status string, leader bool) {
	plr.Status = status

	err := db.UpdatePlayer(plr)
	if err != nil {
		log.Printf("[Error] Failed to update player status in db: %v", err)
		u.SockSend(conn, "Error", plr.PID, plr.GID, "Failed to update player status in db")
	}

	UpdatePlayerList(gid)
}

func Leave(conn *websocket.Conn, plr u.Player, gm u.Game) {
	plr.Status = "Browsing"
	plr.GID = ""

	err := db.UpdatePlayer(plr)
	if err != nil {
		log.Printf("Failed to update player after leaving game: %v", err)
		u.SockSend(conn, "Error", plr.PID, gm.GID, "Could not update player to leave game")
		return
	}

	delete(u.Connections[plr.GID], plr.PID)

	gm.CurPlrs--
	if gm.CurPlrs <= 0 {
		End(conn, gm)
	}

	db.UpdateGame(gm)

	if gm.Leader == plr.PID {
		gm.Leader = ""
		plrs, err := db.GetGamePlayers(gm.GID)
		//TODO find out better way to choose when to end
		if err != nil {
			log.Printf("[Error] Failed to find remaining players in game: %v", err)
		} else {
			if len(plrs) > 0 {
				gm.Leader = plrs[0].PID

				err = db.UpdateGame(gm)
				if err != nil {
					log.Printf("[Error] Failed to update game after finding new leader: %v", err)
					return
				}

				//If there is a new leader, send them an update
				u.SockSend(u.Connections[gm.GID][gm.Leader], "Leader", gm.GID, gm.Leader, "")

				//Send all players the new player list
				UpdatePlayerList(gm.GID)
			} else {
				//TODO this should be redundant
				End(conn, gm)
			}
		}
	}
}
