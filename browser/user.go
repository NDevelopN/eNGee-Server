package browser

import (
	"log"
	"net/http"

	"github.com/google/uuid"

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
	} else {
		// Otherwise, check if player exists
		_, k := u.Plrs[p.PID]
		if !k {
			log.Printf("[Error] Invalid Player ID: %v", p.PID)
			http.Error(w, "Invalid Player ID", http.StatusBadRequest)
			return
		}
	}

	//Add/update player in map
	u.Plrs[p.PID] = p

	//Send user information back to the user, allowing the client to receive the PID
	err = u.PackSend(w, p, "Could not send player update response")
	if err != nil {
		log.Printf("[Error] Failed to send user update information: %v", err)
	}
}
