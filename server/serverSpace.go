package server

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	u "Engee-Server/utils"
)

var defRules = u.Rules{
	Rounds:     1,
	MinPlrs:    2,
	MaxPlrs:    8,
	Timeout:    -1,
	Additional: "",
}

func landing(w http.ResponseWriter, r *http.Request) {
	var p u.Player
	err := u.Extract(r, &p)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		log.Print(err)
		return
	}
	log.Printf("Landing achieved, %v", p.Name)

	if p.PID == "" {
		//Generate UUID for first time player
		p.Status = "New"
		p.PID = uuid.NewString()
	} else {
		_, k := u.Plrs[p.PID]
		if !k {
			//TODO is there anything else to do here?
			log.Printf("Invalid player ID")
			http.Error(w, "Invalid player ID", http.StatusBadRequest)
		}
	}

	u.Plrs[p.PID] = p

	err = u.PackSend(w, p)

	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send reponse: %v\n", err)
	}
}

func browser(w http.ResponseWriter, r *http.Request) {
	var gList u.GameInfo
	var gInfo u.GView

	for i, g := range u.Games {
		gInfo.GID = i
		gInfo.Name = g.Name
		gInfo.Type = g.Type
		gInfo.CurPlrs = len(g.Players)
		gInfo.MaxPlrs = g.Rules.MaxPlrs
		gList.Games = append(gList.Games, gInfo)
	}

	err := u.PackSend(w, gList)
	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send reponse: %v\n", err)
	}
}

func createGame(w http.ResponseWriter, r *http.Request) {
	var g u.Game
	err := u.Extract(r, &g)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		log.Print(err)
		return
	}

	if g.GID == "" {
		//Generate UUID for first time player
		g.GID = uuid.NewString()
		g.Status = "Lobby"
		g.Leader = ""
	}

	u.Games[g.GID] = g
	u.Connections[g.GID] = make(map[string]*websocket.Conn)

	err = u.PackSend(w, u.Games[g.GID])

	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send response: %v", err)
		return
	}
}

func joinGame(w http.ResponseWriter, r *http.Request) {
	var j u.Join
	err := u.Extract(r, &j)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		log.Print(err)
		return
	}

	found, _ := u.CheckForPlayer(j.PID)
	if !found {
		http.Error(w, "Player not found", http.StatusNotFound)
		log.Printf("Could not find player: %v", j.PID)
		return
	}

	found, gm := u.CheckForGame(j.GID)
	if !found {
		http.Error(w, "Game not found", http.StatusNotFound)
		log.Printf("Could not find game: %v", j.GID)
		return
	}

	for _, p := range gm.Players {
		if p.PID == j.PID {
			http.Error(w, "Player alreayd in the game", http.StatusBadRequest)
			log.Printf("Player alreayd in the game: %v", j.PID)
			return
		}
	}

	gm.Players = append(gm.Players, u.Plrs[j.PID])

	if gm.Leader == "" {

		gm.Leader = j.PID
	}

	u.Games[j.GID] = gm

	err = u.PackSend(w, u.ACK{Message: "ACK"})

	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send response: %v", err)
		return
	}
}
