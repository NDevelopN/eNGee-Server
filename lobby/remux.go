package lobby

import (
	u "Engee-Server/utils"
	"net/http"
	"strings"
)

var gmux = map[string]func(http.ResponseWriter, *http.Request){
	"/connect": lobby,
}

//TODO Split handler into game types

func lobby(w http.ResponseWriter, r *http.Request) {
	//TODO Placeholders
	var dud u.MHandler
	var dod ConFunc

	Lobby(w, r, dod, dud)
}

func ReMux(w http.ResponseWriter, r *http.Request) {
	path := strings.Replace(r.URL.Path, "/game", "", 1)
	if handler, ok := gmux[path]; ok {
		handler(w, r)
		return
	}
	http.Error(w, "Invalid route: "+r.URL.Path, http.StatusNotFound)
}
