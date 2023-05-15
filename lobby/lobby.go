package lobby

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	u "Engee-Server/utils"

	"github.com/gorilla/websocket"
)

type ConFunc func(*websocket.Conn, string, string)
type StartFunc func(string)
type GHandler func(u.GameMsg)

func UpdatePlayers(gid string, msg []byte) {
	for _, p := range u.Games[gid].Players {
		u.Connections[gid][p.PID].WriteMessage(websocket.TextMessage, msg)
	}
}

func updateStatus(gm *u.Game, status string) {
	gm.Status = status

	gMsg := u.GameMsg{
		Type:    "Status",
		GID:     gm.GID,
		PID:     "",
		Content: status,
	}

	msg, err := json.Marshal(gMsg)
	if err != nil {
		log.Printf("Cannot marshal game status update message: %v", err)
		return
	}

	UpdatePlayers(gm.GID, msg)
}

func plrListUpdate(gid string) {
	list, err := json.Marshal(u.PlrList{Players: u.Games[gid].Players})
	if err != nil {
		log.Printf("Cannot marshal player list: %v", err)
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
		log.Printf("Cannot marshal player connect update message: %v", err)
	}

	UpdatePlayers(gid, msg)
}

func connect(gid string, pid string, conn *websocket.Conn) {
	gCon := u.Connections[gid]
	gCon[pid] = conn
	u.Connections[gid] = gCon

	lobbyConnect(conn, pid, gid)

	plrListUpdate(gid)
}

func playerStatus(gm *u.Game, pid string, status string, sf StartFunc) {
	ready := 0
	for i, p := range gm.Players {
		if p.PID == pid {
			p.Status = status
			gm.Players[i] = p
			plrListUpdate(gm.GID)
		}

		if p.Status == "Ready" {
			ready++
		}
	}

	//If half of all players are ready, start game
	if ready > (len(gm.Players) / 2) {
		start(gm, sf)

	}
}

func pause(gm u.Game) {
	log.Printf("Status before pause function: %v", gm.Status)
	if gm.Status == "Pause" {
		updateStatus(&gm, gm.OldStatus)
	} else {
		gm.OldStatus = gm.Status
		updateStatus(&gm, "Pause")
	}
	log.Printf("Status after pause function: %v", gm.Status)
	log.Printf("Old Status after pause function: %v", gm.OldStatus)

	u.Games[gm.GID] = gm

}

func start(gm *u.Game, sf StartFunc) {
	//TODO countdown
	if len(gm.Players) >= gm.Rules.MinPlrs {
		updateStatus(gm, "Play")
		sf(gm.GID)
	}
	// TODO else
}

func end(gm *u.Game) {
	restart(gm)
	//TODO change to delete the game
}

func restart(gm *u.Game) {
	updateStatus(gm, "Restart")
	for i := range gm.Players {
		gm.Players[i].Status = "New"
	}
	time.Sleep(1 * time.Second)
	updateStatus(gm, "Lobby")
}

func rules(gid string, content string) {
	var rules u.Game
	err := json.Unmarshal([]byte(content), &rules)
	if err != nil {
		log.Printf("Error unmarshalling rules update: %v", err)
		return
	}

	//TODO rules validation

	gm := u.Games[gid]
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
		log.Print("No changes to game rules provided")
		return
	}

	updateStatus(&gm, "Restart")
	for i := range gm.Players {
		gm.Players[i].Status = "New"
	}

	gUpd, err := json.Marshal(gm)
	if err != nil {
		log.Printf("Could not marshal games for rules update: %v", err)
		return
	}

	gMsg := u.GameMsg{
		Type:    "Update",
		PID:     "",
		GID:     gid,
		Content: string(gUpd),
	}

	msg, err := json.Marshal(gMsg)
	if err != nil {
		log.Printf("Could not marshal rules update message: %v", err)
		return
	}

	//TODO clean this up
	UpdatePlayers(gid, msg)

	time.Sleep(1 * time.Second)

	updateStatus(&gm, "Lobby")

	u.Games[gm.GID] = gm
}

// TO send update only to one Player
func SingleWrite(t string, pid string, gid string, content string) {
	msg := u.GameMsg{
		Type:    t,
		PID:     pid,
		GID:     gid,
		Content: content,
	}

	enc, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Cannot marshal leader message: %v", err)
		return
	}

	log.Printf("Got this far: %v: %v", gid, pid)

	u.Connections[gid][pid].WriteMessage(websocket.TextMessage, enc)
}

func lobbyConnect(conn *websocket.Conn, pid string, gid string) {
	gm := u.Games[gid]

	//Make first connector be leader
	if gm.Leader == "" {
		gm.Leader = pid
	}

	info, err := json.Marshal(gm)
	if err != nil {
		log.Printf("Cannot marshal game info: %v", err)
		return
	}

	SingleWrite("Info", pid, gid, string(info))
}

func Lobby(w http.ResponseWriter, r *http.Request, gameConnect ConFunc, sf StartFunc, gameHandler GHandler) {

	var lobbyHandler u.MHandler = func(conn *websocket.Conn, data []byte) {
		var msg u.GameMsg
		err := json.Unmarshal(data, &msg)
		if err != nil {
			log.Printf("Cannot unmarshal message: %v", err)
			return
		}

		found, _ := u.CheckForPlayer(msg.PID)
		if !found {
			log.Printf("Player not found: %v", msg.PID)
			return
		}

		found, gm := u.CheckForGame(msg.GID)
		if !found {
			log.Printf("Game not found: %v", msg.GID)
			return
		}

		//Check if sender is the leader of the game
		leader := (msg.PID == gm.Leader)

		switch msg.Type {
		case "Connect":
			connect(msg.GID, msg.PID, conn)
		case "Status":
			playerStatus(&gm, msg.PID, msg.Content, sf)
			u.Games[msg.GID] = gm
		case "Leave":
			//Remove the player, if they were not the last update the others
			if u.RemovePlayer(msg.GID, msg.PID) {
				log.Printf("REmove player %v all done", msg.PID)
				log.Printf("NEw leader: %v", u.Games[msg.GID].Leader)
				SingleWrite("Leader", u.Games[msg.GID].Leader, msg.GID, "")
				plrListUpdate(msg.GID)
			}
		case "Pause":
			if !leader {
				log.Printf("%v, is not leader", msg.PID)
				return
			}
			pause(gm)
		case "Start":
			if !leader {
				log.Printf("%v, is not leader", msg.PID)
				return
			}

			start(&gm, sf)
			u.Games[msg.GID] = gm
		case "End":
			if !leader {
				log.Printf("%v, is not leader", msg.PID)
				return
			}

			end(&gm)
			u.Games[msg.GID] = gm
		case "Restart":
			if !leader {
				log.Printf("%v, is not leader", msg.PID)
				return
			}

			restart(&gm)

			u.Games[msg.GID] = gm
		case "Rules":
			if !leader {
				log.Printf("%v, is not leader", msg.PID)
				return
			}

			rules(msg.GID, msg.Content)
		case "Remove":
			//TODO
		case "Delete":
			//TODO

		default:
			gameHandler(msg)
			return
		}
	}

	u.Sock(w, r, lobbyHandler)
}
