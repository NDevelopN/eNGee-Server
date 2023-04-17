package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

var mux = map[string]func(http.ResponseWriter, *http.Request){
	"/":              landing,
	"/server":        landing,
	"/serverbrowser": browser,
	"/server/create": createGame,
	"/game/join":     joinGame,
}

func landing(w http.ResponseWriter, r *http.Request) {
	var pInfo PlayerInfo
	err := extract(r, &pInfo)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		log.Print(err)
		return
	}

	var Plr Player
	id := pInfo.ID
	if id == "" {
		//Generate UUID for first time player
		Plr.Status = "New"
		Plr.Games = make(map[string]string)
		id = uuid.NewString()

	} else {
		_, k := Plrs[id]
		if !k {
			//TODO is there anything else to do here?
			log.Printf("Invalid player ID")
			http.Error(w, "Invalid player ID", http.StatusBadRequest)
		}
		Plr = Plrs[id]
	}

	Plr.Name = pInfo.Name
	Plrs[id] = Plr

	pInfo.ID = id

	err = packSend(w, pInfo)

	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send reponse: %v\n", err)
	}
}

func browser(w http.ResponseWriter, r *http.Request) {
	var pInfo PlayerInfo
	err := extract(r, &pInfo)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		log.Print(err)
		return
	}

	//TODO: add some importance to games player is already in
	var gList GameList
	var gInfo GameInfo
	for i, g := range Games {
		gInfo.ID = i
		gInfo.Name = g.Name
		gInfo.GameType = g.GameType
		gInfo.Status = g.Status
		gList.Games = append(gList.Games, gInfo)
	}

	err = packSend(w, gList)
	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send reponse: %v\n", err)
	}
}

func createGame(w http.ResponseWriter, r *http.Request) {
	var rules GameRules
	err := extract(r, &rules)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		log.Print(err)
		return
	}

	var gm Game
	id := rules.ID
	if id == "" {
		//Generate UUID for first time player
		gm.Status = "New"
		id = uuid.NewString()

	} else {
		_, k := Games[id]
		if !k {
			//TODO is there anything else to do here?
			http.Error(w, "Invalid Game ID", http.StatusBadRequest)
			log.Printf("Invalid Game ID")
		}
		gm = Games[id]
	}

	gm.Name = rules.Name
	gm.GameType = rules.GameType
	gm.MinPlayers = rules.MinPlayers
	gm.MaxPlayers = rules.MaxPlayers

	Games[id] = gm

	var gInfo GameInfo
	gInfo.ID = id
	gInfo.Name = rules.Name
	gInfo.GameType = rules.GameType
	gInfo.Status = gm.Status
	gInfo.PlayerCount = 0

	err = packSend(w, gInfo)
	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send response: %v", err)
		return
	}

}

func joinGame(w http.ResponseWriter, r *http.Request) {
	var join Join
	err := extract(r, &join)
	if err != nil {
		http.Error(w, "Could not read join from request", http.StatusBadRequest)
		log.Printf("Could not read from request %v", err)
	}

	p, k := Plrs[join.PID]
	if !k {
		http.Error(w, "Player not found", http.StatusNotFound)
		log.Print("Player not found")

		return
	}

	_, k = p.Games[join.GID]
	if k {
		http.Error(w, "Player already in game", http.StatusBadRequest)
		log.Print("Player already in game")
		return
	}

	g, k := Games[join.GID]
	if !k {
		http.Error(w, "Game not found", http.StatusNotFound)
		log.Print("Game not found")
		return
	}

	p.Games[join.GID] = g.Name
	g.Players = append(g.Players, p.Name)
	g.PlayerCount++

	Plrs[join.PID] = p
	Games[join.GID] = g

	var game GameInfo
	game.ID = join.GID
	game.Name = g.Name
	game.GameType = g.GameType
	game.Status = g.Status
	game.PlayerCount = g.PlayerCount

	err = packSend(w, game)
	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send response: %v", err)
		return
	}
}

func extract[O *PlayerInfo | *GameRules | *Join](r *http.Request, obj O) error {
	reqBody, err := io.ReadAll(r.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(reqBody, obj)

	if err != nil {
		return err
	}

	return nil
}

func packSend[O PlayerInfo | GameInfo | GameList](w http.ResponseWriter, msg O) error {
	response, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		return err
	}

	return nil

}
