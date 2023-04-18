package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"

	u "Engee-Server/utils"
)

var smux = map[string]func(http.ResponseWriter, *http.Request){
	"/":        landing,
	"/browser": browser,
	"/create":  createGame,
}

func ReMux(w http.ResponseWriter, r *http.Request) {
	path := strings.Replace(r.URL.Path, "/server", "", 1)

	if handler, ok := smux[path]; ok {
		handler(w, r)
		return
	}
	http.Error(w, "Invalid route: "+r.URL.Path, http.StatusNotFound)

}

func landing(w http.ResponseWriter, r *http.Request) {
	var pInfo u.PlayerInfo
	err := u.Extract(r, &pInfo)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		log.Print(err)
		return
	}

	var Plr u.Player
	id := pInfo.ID
	if id == "" {
		//Generate UUID for first time player
		Plr.Status = "New"
		Plr.Games = make(map[string]string)
		id = uuid.NewString()

	} else {
		_, k := u.Plrs[id]
		if !k {
			//TODO is there anything else to do here?
			log.Printf("Invalid player ID")
			http.Error(w, "Invalid player ID", http.StatusBadRequest)
		}
		Plr = u.Plrs[id]
	}

	Plr.Name = pInfo.Name
	u.Plrs[id] = Plr

	pInfo.ID = id

	err = u.PackSend(w, pInfo)

	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send reponse: %v\n", err)
	}
}

func browser(w http.ResponseWriter, r *http.Request) {
	var pInfo u.PlayerInfo
	err := u.Extract(r, &pInfo)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		log.Print(err)
		return
	}

	//TODO: add some importance to games player is already in
	var gList u.GameList
	var gInfo u.GameInfo
	for i, g := range u.Games {
		gInfo.ID = i
		gInfo.Name = g.Name
		gInfo.GameType = g.GameType
		gInfo.Status = g.Status
		gList.Games = append(gList.Games, gInfo)
	}

	err = u.PackSend(w, gList)
	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send reponse: %v\n", err)
	}
}

func createGame(w http.ResponseWriter, r *http.Request) {
	var rules u.GameRules
	err := u.Extract(r, &rules)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		log.Print(err)
		return
	}

	var gm u.Game
	id := rules.ID
	if id == "" {
		//Generate UUID for first time player
		gm.Status = "Ready"
		id = uuid.NewString()
		gm.Players = make(map[string]string)

	} else {
		_, k := u.Games[id]
		if !k {
			//TODO is there anything else to do here?
			http.Error(w, "Invalid Game ID", http.StatusBadRequest)
			log.Printf("Invalid Game ID")
		}
		gm = u.Games[id]
	}

	gm.Name = rules.Name
	gm.GameType = rules.GameType
	gm.MinPlayers = rules.MinPlayers
	gm.MaxPlayers = rules.MaxPlayers

	u.Games[id] = gm

	var gInfo u.GameInfo
	gInfo.ID = id
	gInfo.Name = rules.Name
	gInfo.GameType = rules.GameType
	gInfo.Status = gm.Status
	gInfo.PlayerCount = 0

	err = u.PackSend(w, gInfo)
	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send response: %v", err)
		return
	}

}
