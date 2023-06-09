package browser

import (
	"log"
	"net/http"

	db "Engee-Server/database"
	u "Engee-Server/utils"
)

/**
 *
 * The Browser function provides callers with a list of information about currently available games
 * This function sends a list of GView objects, which contain less information than full Game objects
 * The function makes no distinction between an empty and populated list
 * TODO: Should this function be provided with restrictions on the games the user requires?
 *
 */
func Browser(w http.ResponseWriter, r *http.Request) {

	games, err := db.GetAllGames()
	if err != nil {
		log.Printf("[Error] Failed to get games from database: %v", err)
		return
	}

	err = u.SendGameInfo(w, games, "Could not send game browser info")
	if err != nil {
		log.Printf("[Error] Failed to send game browser information: %v", err)
		return
	}
}

/**
 *
 * This JoinGame function handles users requests to join a game
 * This function checks if the user and the game exists, returning an error if either is not found
 * This function then checks if the user is already in the game, returning an error if true
 * This funciton then checks if there is enough room for the user to join, returning an error if there is not
 * This function then adds the user to the game's map of players
 * This function sends the user an ACK message to confirm the join request was successful
 * TODO: Do all of the above need errors, or would another message be better?
 *
 */
func JoinGame(w http.ResponseWriter, r *http.Request) {
	var j u.Join
	err := u.Extract(r, &j)
	if err != nil {
		log.Printf("[Error] Failed to read join request: %v", err)
		http.Error(w, "Failed to read join request", http.StatusBadRequest)
		return
	}

	//Check if the player exists
	plr, err := db.GetPlayer(j.PID)
	if err != nil {
		log.Printf("[Error] Failed to find player in DB: %v", err)
		http.Error(w, "Failed to find player in the Database", http.StatusInternalServerError)
		return
	}

	//Check if the game exists
	gm, err := db.GetGame(j.GID)
	if err != nil {
		log.Printf("[Error] Failed to find game in DB: %v", err)
		http.Error(w, "Failed to find game in the Database", http.StatusInternalServerError)
		return
	}

	//Check if the player is in the game
	if plr.GID == gm.GID {
		log.Printf("[Error] Player is already in game")
		http.Error(w, "Player is already in selected game", http.StatusBadRequest)
		return
	}

	//Check if there is enough room for the player to join
	pCount := db.GetGamePCount(j.GID)
	if pCount >= gm.MaxPlrs {
		log.Printf("[Error] No space for player to join: %v", err)
		http.Error(w, "No room for player to join", http.StatusConflict)
		return
	}

	//Update the player to be part of the game
	plr.GID = gm.GID
	plr.Status = "Joining"

	//If there hasn't yet been an assigned leader, set the player as the game leader
	if gm.Leader == "" {
		gm.Leader = plr.PID
	}

	gm.CurPlrs = pCount + 1

	db.UpdateGame(gm)
	db.UpdatePlayer(plr)

	//Send an ACK message to the user to complete the process
	err = u.PackSend(w, u.ACK{Message: "ACK"}, "Could not send join Acknowledgement")
	if err != nil {
		log.Printf("[Error] Failed to send game join ACK to player: %v", err)
	}
}
