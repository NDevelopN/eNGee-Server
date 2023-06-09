package browser

import (
	"log"
	"net/http"

	"github.com/google/uuid"

	db "Engee-Server/database"
	u "Engee-Server/utils"
)

/**
 *
 * The EditUser function handles user creation and information updates.
 * If the request does not contain a PID, this function generates a new one using UUID.
 * If the request does contain a PID, this function searches for a corresponding entry in the map.
 * The Player information is added/updated and then returned to the caller.
 * In any case of failure, an error is instead returned.
 *
 */
func EditUser(w http.ResponseWriter, r *http.Request) {
	var p u.Player
	err := u.Extract(r, &p)
	if err != nil {
		log.Printf("[Error] Failed to read user request: %v", err)
		return
	}

	if p.PID == "" {
		// Create UUID for new player
		p.Status = "New"
		p.PID = uuid.NewString()
		err := db.CreatePlayer(p)
		if err != nil {
			log.Printf("[Error] Could not create user in database: %v", err)
			http.Error(w, "Could not create new user in database", http.StatusInternalServerError)
			return
		}
	} else {
		_, err = db.GetPlayer(p.PID)
		if err != nil {
			log.Printf("[Error] Could not retrieve player from database: %v", err)
			http.Error(w, "Could not retrieve player from database", http.StatusInternalServerError)
			return
		}

		err = db.UpdatePlayer(p)
		if err != nil {
			log.Printf("[Error] Could not update player in database: %v", err)
			http.Error(w, "Could not update player in database", http.StatusInternalServerError)
			return
		}
	}

	//Send user information back to the user, allowing the client to receive the PID
	err = u.PackSend(w, p, "Could not send player update response")
	if err != nil {
		log.Printf("[Error] Failed to send user update information: %v", err)
	}
}
