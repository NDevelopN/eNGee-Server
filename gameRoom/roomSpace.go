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
	"/leave":   leave,
	"/ready":   ready,
	"/start":   start,
	"/delete":  dlt,
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
	g.Players[join.PID] = "Joined"
	g.PlayerCount++

	u.Plrs[join.PID] = p
	u.Games[join.GID] = g

	var game u.GameInfo
	game.ID = join.GID
	game.Name = g.Name
	game.GameType = g.GameType
	game.Status = g.Status
	game.PlayerCount = g.PlayerCount

	log.Print(p)
	log.Print(g)

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
		g.Players[ready.PID] = "Not Ready"
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

	p, fail := checkPlayerList(w, restart)
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

func pause(w http.ResponseWriter, r *http.Request) {
	var pause u.GameOp
	err := u.Extract(r, &pause)
	if err != nil {
		http.Error(w, "Could not read restart from request", http.StatusInternalServerError)
		log.Printf("Could not read from request %v", err)
	}

	p, fail := checkPlayerList(w, pause)
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
		g.Status = "InGame"
	} else {
		g.Status = "Pause"
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

func end(w http.ResponseWriter, r *http.Request) {
	var end u.GameOp
	err := u.Extract(r, &end)
	if err != nil {
		http.Error(w, "Could not read restart from request", http.StatusInternalServerError)
		log.Printf("Could not read from request %v", err)
	}

	p, fail := checkPlayerList(w, end)
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
	g.Status = "Ready"

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
