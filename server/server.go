package server

import (
	"log"
	"net/http"
	"strings"
	"time"

	b "Engee-Server/browser"
	db "Engee-Server/database"
	consequences "Engee-Server/games/consequences"
	g "Engee-Server/gamespace"
	u "Engee-Server/utils"
)

var gHandlers = map[string]u.GHandler{
	"consequences": consequences.HandleInput,
}

type myHandler struct{}

/**TODO: Review and see which are needed */
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

var mux = map[string]func(http.ResponseWriter, *http.Request){
	"":       ServerRouting,
	"server": ServerRouting,
	"game":   GameRouting,
}

var smux = map[string]func(http.ResponseWriter, *http.Request){
	"/":        b.EditUser,
	"/browser": b.Browser,
	"/create":  b.EditGame,
	"/join":    b.JoinGame,
}

func ServerRouting(w http.ResponseWriter, r *http.Request) {
	path := strings.Replace(r.URL.Path, "/server", "", 1)

	if handler, ok := smux[path]; ok {
		handler(w, r)
		return
	}
	http.Error(w, "Invalid route: "+r.URL.Path, http.StatusNotFound)
}

func GameRouting(w http.ResponseWriter, r *http.Request) {
	gid := strings.Replace(r.URL.Path, "/game/", "", 1)
	gm, err := db.GetGame(gid)
	if err != nil {
		log.Printf("[Error] Game in path not found: %v", err)
		http.Error(w, "Game in path not found", http.StatusNotFound)
		return
	}
	t := strings.ToLower(gm.Type)

	if handler, ok := gHandlers[t]; ok {
		g.GameSpace(w, r, handler)
	}

	http.Error(w, "Invalid route: "+r.URL.Path, http.StatusNotFound)
}

func (h *myHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// Handle CORS issues TODO: Review CORS in detail
	enableCors(&writer)
	if request.Method == "OPTIONS" {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(""))
		return
	}

	// Get the base path and pass to relevant ReMux
	path := strings.Split(request.URL.Path, "/")[1]

	//Implement route forwarding, ensure there is a route established for the request
	if handler, ok := mux[path]; ok {
		handler(writer, request)
		return
	}
	http.Error(writer, "Invalid route: "+request.URL.Path, http.StatusNotFound)

}

func Serve() error {
	db.InitDB()
	server := http.Server{
		Addr:        ":8090",
		Handler:     &myHandler{},
		ReadTimeout: 5 * time.Second,
	}
	return server.ListenAndServe()
}
