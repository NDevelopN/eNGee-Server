package GameRoom

import (
	"log"
	"net/http"
	"strings"

	//s "Engee-Server/server"
	u "Engee-Server/utils"
)

var gmux = map[string]func(http.ResponseWriter, *http.Request){
	"/join":    join,
	"/status":  status,
	"/leave":   leave,
	"/ready":   ready,
	"/start":   start,
	"/delete":  dlt,
	"/rules":   rules,
	"/update":  update,
	"/remove":  remove,
	"/restart": restart,
	"/pause":   pause,
	"/end":     end,
}

func ReMux(w http.ResponseWriter, r *http.Request) {
	path := strings.Replace(r.URL.Path, "/game", "", 1)
	log.Print(path)

	if handler, ok := gmux[path]; ok {
		handler(w, r)
		return
	}
	http.Error(w, "Invalid route: "+r.URL.Path, http.StatusNotFound)
}

func checkPlayerList(w http.ResponseWriter, op u.GameOp) (u.Player, bool) {

	var b bool = false

	p, k := u.Plrs[op.PID]
	if !k {
		http.Error(w, "Player not found", http.StatusNotFound)
		log.Print("Player not found")

		b = true
	}

	return p, b
}

func checkGameList(w http.ResponseWriter, op u.GameOp) (u.Game, bool) {

	var b bool = false

	g, k := u.Games[op.GID]
	if !k {
		http.Error(w, "Game not found", http.StatusNotFound)
		log.Print("Game not found")
		b = true
	}

	return g, b
}

func status(w http.ResponseWriter, r *http.Request) {
	var req u.GameOp
	err := u.Extract(r, &req)
	if err != nil {
		http.Error(w, "Could not read request for players", http.StatusInternalServerError)
		log.Printf("Could not read from request: %v", err)
	}

	g, fail := checkGameList(w, req)
	if fail {
		return
	}

	var plr u.Player
	var plrStat u.PlayerStatus

	var status u.GameStatus

	status.Status = g.Status
	status.Leader = g.Leader

	for i, p := range g.Players {
		plr = u.Plrs[i]
		plrStat.Name = plr.Name
		plrStat.ID = i
		plrStat.Status = p
		status.Players = append(status.Players, plrStat)
	}

	err = u.PackSend(w, status)
	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send response: %v", err)
		return
	}
	log.Println("Function complete")
}
func join(w http.ResponseWriter, r *http.Request) {
	var join u.GameOp
	err := u.Extract(r, &join)
	if err != nil {
		http.Error(w, "Could not read join from request", http.StatusBadRequest)
		log.Printf("Could not read from request %v", err)
	}

	p, fail := checkPlayerList(w, join)

	if fail {
		return
	}

	g, fail := checkGameList(w, join)
	if fail {
		return
	}

	//Check to see if player already has joined this game
	_, k := p.Games[join.GID]
	if k {
		http.Error(w, "Player already in chosen game", http.StatusBadRequest)
		log.Print("Player already in chosen game", http.StatusBadRequest)
		return
	}

	if g.Leader == "" {
		g.Leader = join.PID
	}

	p.Games[join.GID] = g.Name
	g.Players[join.PID] = "Lobby"
	g.PlayerCount++

	u.Plrs[join.PID] = p
	u.Games[join.GID] = g

	var game u.GameInfo
	game.ID = join.GID
	game.Name = g.Name
	game.GameType = g.GameType
	game.Status = g.Status
	game.PlayerCount = g.PlayerCount

	err = u.PackSend(w, game)
	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send response: %v", err)
		return
	}
}

func leave(w http.ResponseWriter, r *http.Request) {
	var leave u.GameOp
	err := u.Extract(r, &leave)
	if err != nil {
		http.Error(w, "Could not read leave from request", http.StatusInternalServerError)
		log.Printf("Could not read from request %v", err)
	}

	p, fail := checkPlayerList(w, leave)
	if fail {
		return
	}

	g, fail := checkGameList(w, leave)
	if fail {
		return
	}

	_, k := p.Games[leave.GID]
	if !k {
		http.Error(w, "Player not in chosen game", http.StatusBadRequest)
		log.Print("Player not in chosen game", http.StatusBadRequest)
		return
	}

	delete(p.Games, leave.GID)

	if g.Leader == leave.PID {
		g.Leader = ""
		for l := range g.Players {
			if k { //Assign a random player to the leader
				g.Leader = l
			} else { //Or remove the game if there are no more players
				delete(u.Games, leave.GID)
			}
			break
		}
	}

	g.Players[leave.PID] = "Left"
	log.Print(p)
	log.Print(g)

	//TODO change to ACK
	err = u.PackSend(w, "")
	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send response: %v", err)
	}

}

func ready(w http.ResponseWriter, r *http.Request) {
	var ready u.GameOp
	err := u.Extract(r, &ready)
	if err != nil {
		http.Error(w, "Could not read ready from request", http.StatusInternalServerError)
		log.Printf("Could not read from request %v", err)
	}

	p, fail := checkPlayerList(w, ready)
	if fail {
		return
	}

	g, fail := checkGameList(w, ready)
	if fail {
		return
	}

	_, k := p.Games[ready.GID]
	if !k {
		http.Error(w, "Player not in chosen game", http.StatusBadRequest)
		log.Print("Player not in chosen game", http.StatusBadRequest)
		return
	}

	//TODO maybe more conditions
	if g.Players[ready.PID] == "Ready" {
		g.Players[ready.PID] = g.Status
	} else {
		g.Players[ready.PID] = "Ready"
	}

	log.Print(p)
	log.Print(g)

	//TODO remove player from game
	//TODO change to ACK
	err = u.PackSend(w, "")
	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send response: %v", err)
	}

}

func start(w http.ResponseWriter, r *http.Request) {
	var start u.GameOp
	err := u.Extract(r, &start)
	if err != nil {
		http.Error(w, "Could not read start from request", http.StatusInternalServerError)
		log.Printf("Could not read from request %v", err)
	}

	p, fail := checkPlayerList(w, start)
	if fail {
		return
	}

	g, fail := checkGameList(w, start)
	if fail {
		return
	}

	lead := g.Leader
	if start.PID != lead {
		http.Error(w, "You are not the game leader", http.StatusBadRequest)
		log.Printf("Not game leader, refused")
		return
	}

	var readyCount int = 0
	//TODO check if enough players are ready
	for _, plr := range g.Players {
		if plr == "Ready" {
			readyCount++
		}
	}

	if readyCount < g.PlayerCount/2 {
		http.Error(w, "There are not enough ready players", http.StatusInternalServerError)
		log.Printf("Not enough ready players")
		return
	}

	g.Status = "InGame"

	for i := range g.Players {
		g.Players[i] = "InGame"
	}
	//TODO now do sometihng for game specifics

	log.Print(p)
	log.Print(g)

	//TODO remove player from game
	//TODO change to ACK
	err = u.PackSend(w, "")
	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send response: %v", err)
	}
}

func dlt(w http.ResponseWriter, r *http.Request) {
	var dlt u.GameOp
	err := u.Extract(r, &dlt)
	if err != nil {
		http.Error(w, "Could not read delete from request", http.StatusInternalServerError)
		log.Printf("Could not read from request %v", err)
	}

	p, fail := checkPlayerList(w, dlt)
	if fail {
		return
	}

	g, fail := checkGameList(w, dlt)
	if fail {
		return
	}

	lead := g.Leader
	if dlt.PID != lead {
		http.Error(w, "You are not the game leader", http.StatusBadRequest)
		log.Printf("Not game leader, refused")
		return
	}

	delete(u.Games, dlt.GID)

	g.Status = "Ended"
	//TODO now do sometihng for game specifics

	log.Print(p)
	log.Print(g)

	//TODO remove player from game
	//TODO change to ACK
	err = u.PackSend(w, "")
	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send response: %v", err)
	}

}

func rules(w http.ResponseWriter, r *http.Request) {
	var rules u.GameOp
	err := u.Extract(r, &rules)
	if err != nil {
		http.Error(w, "Could not read rules request", http.StatusBadRequest)
		log.Printf("Could not read rules request: %v", err)
	}

	g, fail := checkGameList(w, rules)
	if fail {
		return
	}

	if g.Leader != rules.PID {
		http.Error(w, "Player is not leader of game", http.StatusBadRequest)
		log.Printf("Player (%v) not leader of game (%v)", rules.PID, g.Leader)
		return
	}

	var ri u.GameRules
	ri.ID = rules.PID
	ri.Name = g.Name
	ri.GameType = g.GameType
	ri.MinPlayers = g.MinPlayers
	ri.MaxPlayers = g.MaxPlayers
	//TODO
	ri.Additional = ""

	err = u.PackSend(w, ri)
	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send response: %v", err)
	}
}

func update(w http.ResponseWriter, r *http.Request) {
	var rules u.GameRules
	err := u.Extract(r, &rules)
	if err != nil {
		http.Error(w, "Could not read rules from request", http.StatusBadRequest)
		log.Printf("Could no ready rules from request: %v", err)
	}

	var op u.GameOp
	op.GID = rules.ID

	//TODO change GameRules to include pid
	g, fail := checkGameList(w, op)
	if fail {
		return
	}

	g.Name = rules.Name
	g.GameType = rules.GameType
	g.MinPlayers = rules.MinPlayers
	g.MaxPlayers = rules.MaxPlayers
	//TODO additional

	//TODO change to ACK
	err = u.PackSend(w, "")
	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send response: %v", err)
	}
}

func remove(w http.ResponseWriter, r *http.Request) {
	var remove u.RemovePlr
	err := u.Extract(r, &remove)
	if err != nil {
		http.Error(w, "Could not read remove from request", http.StatusInternalServerError)
		log.Printf("Could not read from request %v", err)
	}

	var op u.GameOp
	op.PID = remove.PlrID
	op.GID = remove.GID

	p, fail := checkPlayerList(w, op)
	if fail {
		return
	}

	g, fail := checkGameList(w, op)
	if fail {
		return
	}

	lead := g.Leader
	if remove.AdminID != lead {
		http.Error(w, "You are not the game leader", http.StatusBadRequest)
		log.Printf("Not game leader, refused")
		return
	}

	delete(p.Games, remove.GID)
	delete(g.Players, remove.PlrID)

	//TODO now do sometihng for game specifics

	log.Print(p)
	log.Print(g)

	//TODO remove player from game
	//TODO change to ACK
	err = u.PackSend(w, "")
	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send response: %v", err)
	}
}

func restart(w http.ResponseWriter, r *http.Request) {
	var restart u.GameOp
	err := u.Extract(r, &restart)
	if err != nil {
		http.Error(w, "Could not read restart from request", http.StatusInternalServerError)
		log.Printf("Could not read from request %v", err)
	}

	_, fail := checkPlayerList(w, restart)
	if fail {
		return
	}

	g, fail := checkGameList(w, restart)
	if fail {
		return
	}

	lead := g.Leader
	if restart.PID != lead {
		http.Error(w, "You are not the game leader", http.StatusBadRequest)
		log.Printf("Not game leader, refused")
		return
	}

	/**
	var readyCount int = 0
	for _, plr := range g.Players {
		if plr == "Ready" {
			readyCount++
		}
	}

	if readyCount < g.PlayerCount/2 {
		http.Error(w, "There are not enough ready players", http.StatusInternalServerError)
		log.Printf("Not enough ready players")
		return
	}*/

	g.Status = "Lobby"
	for i := range g.Players {
		g.Players[i] = "Lobby"
	}
	//TODO now do sometihng for game specifics

	u.Games[restart.GID] = g

	//TODO remove player from game
	//TODO change to ACK
	err = u.PackSend(w, "")
	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send response: %v", err)
	}

}

func pause(w http.ResponseWriter, r *http.Request) {
	var pause u.GameOp
	err := u.Extract(r, &pause)
	if err != nil {
		http.Error(w, "Could not read restart from request", http.StatusInternalServerError)
		log.Printf("Could not read from request %v", err)
	}

	_, fail := checkPlayerList(w, pause)
	if fail {
		return
	}

	g, fail := checkGameList(w, pause)
	if fail {
		return
	}

	lead := g.Leader
	if pause.PID != lead {
		http.Error(w, "You are not the game leader", http.StatusBadRequest)
		log.Printf("Not game leader, refused")
		return
	}

	// TODO more conditions?
	if g.Status == "Pause" {
		g.Status = g.OldStatus
	} else {
		g.OldStatus = g.Status
		g.Status = "Pause"
	}

	u.Games[pause.GID] = g

	//TODO remove player from game
	//TODO change to ACK
	err = u.PackSend(w, "")
	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send response: %v", err)
	}
}

func end(w http.ResponseWriter, r *http.Request) {
	var end u.GameOp
	err := u.Extract(r, &end)
	if err != nil {
		http.Error(w, "Could not read restart from request", http.StatusInternalServerError)
		log.Printf("Could not read from request %v", err)
	}

	_, fail := checkPlayerList(w, end)
	if fail {
		return
	}

	g, fail := checkGameList(w, end)
	if fail {
		return
	}

	lead := g.Leader
	if end.PID != lead {
		http.Error(w, "You are not the game leader", http.StatusBadRequest)
		log.Printf("Not game leader, refused")
		return
	}

	// TODO more conditions?
	g.Status = "Lobby"

	u.Games[end.GID] = g

	//TODO remove player from game
	//TODO change to ACK
	err = u.PackSend(w, "")
	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send response: %v", err)
	}
}
