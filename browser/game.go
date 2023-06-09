package browser

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	db "Engee-Server/database"
	u "Engee-Server/utils"
)

/**
 *
 * The EditGame function handles game creation and modification.
 * If the request does not contain a GID, this function generates a new one using UUID.
 * If the request does contain a GID, this function searches for a corresponding entry in the map.
 * The Game information is added/updated and then returned to the caller.
 * In any case of failure, an eror is instead returned.
 *
 */
func EditGame(w http.ResponseWriter, r *http.Request) {
	var g u.Game
	err := u.Extract(r, &g)
	if err != nil {
		log.Printf("[Error] Failed to read game creation request: %v", err)
		return
	}

	if g.GID == "" {
		//Generate UUID for new games
		g.GID = uuid.NewString()
		g.Status = "Lobby"
		g.Leader = ""
		err := db.CreateGame(g)
		if err != nil {
			log.Printf("[Error] Could not create game in database: %v", err)
			http.Error(w, "Could not create new game in database", http.StatusInternalServerError)
			return
		}
	} else {
		_, err := db.GetGame(g.GID)
		if err != nil {
			log.Printf("[Error] Could not retieve  game from database: %v", err)
			http.Error(w, "Could not retrieve game from database", http.StatusInternalServerError)
			return
		}

		err = db.UpdateGame(g)
		if err != nil {
			log.Printf("[Error] Could not update  game in database: %v", err)
			http.Error(w, "Could not update game in database", http.StatusInternalServerError)
			return
		}
	}

	//Make new map for player connections
	u.Connections[g.GID] = make(map[string]*websocket.Conn)

	//Send game information back to the user, allowing the client to receive the GID
	err = u.PackSend(w, g, "Could not send game update response")
	if err != nil {
		log.Printf("[Error] Failed to send game update information: %v", err)
	}
}
