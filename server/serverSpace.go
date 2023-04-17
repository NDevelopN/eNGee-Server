package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

var mux = map[string]func(http.ResponseWriter, *http.Request){
	"/": landing,
}

func landing(w http.ResponseWriter, r *http.Request) {
	var pInfo PlayerInfo
	reqBody, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		log.Print(err)
		return
	}

	err = json.Unmarshal(reqBody, &pInfo)

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
	Plrs[pInfo.ID] = Plr

	pInfo.ID = id

	response, err := json.Marshal(pInfo)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		log.Printf("Failed to marshal reponse: %s\n", err)
	}

	w.WriteHeader(http.StatusOK)
	i, err := w.Write(response)
	if err != nil {
		http.Error(w, "Failed in response", http.StatusInternalServerError)
		log.Printf("Failed to send reponse: %d\nError: %s\n", i, err)
	}
}
