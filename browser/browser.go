package browser

import (
	"log"
	"net/http"

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
	var gList u.GameInfo
	var gInfo u.GView
	//TODO do I need gList to be a struct

	//For each game, put its information into the GameInfo list
	for i, g := range u.Games {
		gInfo.GID = i
		gInfo.Name = g.Name
		gInfo.Status = g.Status
		gInfo.Type = g.Type
		gInfo.CurPlrs = len(g.Players)
		gInfo.MaxPlrs = g.Rules.MaxPlrs
		gList.Games = append(gList.Games, gInfo)
	}

	err := u.PackSend(w, gList, "Could not send game browser information")
	if err != nil {
		log.Printf("[Error] Failed to send game browser information: %v", err)
		return
	}
}

/**
 *
 * This JoinGame fucntion handles users requests to join a game
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
		return
	}

	//Check if the player exists
	found, p := u.CheckForPlayer(j.PID)
	if !found {
		http.Error(w, "Player was not found", http.StatusNotFound)
		return
	}

	//Check if the game exists
	found, gm := u.CheckForGame(j.GID)
	if !found {
		http.Error(w, "Game was not found", http.StatusNotFound)
		return
	}

	//Check if the player is already in the game
	found = u.CheckGameForPlayer(gm, j.PID)
	if found {
		http.Error(w, "Player is already in the game", http.StatusBadRequest)
		log.Printf("[Error] Player is trying to join a game they are already in")
		return
	}

	//Check if there is enough space for the player to join
	if len(gm.Players) == gm.Rules.MaxPlrs {
		http.Error(w, "Cannot join game, it is already full", http.StatusConflict)
		log.Printf("[Error] Player trying to join a full game")
		return
	}

	//Add the player to the game's map
	gm.Players = append(gm.Players, p)

	//If there hasn't yet been an assigned leader, set the player as the game leader
	if gm.Leader == "" {
		gm.Leader = p.PID
	}

	//Update the game map with the update game
	u.Games[j.GID] = gm

	//Send an ACK message to the user to complete the process
	err = u.PackSend(w, u.ACK{Message: "ACK"}, "Could not send join Acknowledgement")
	if err != nil {
		log.Printf("[Error] Failed to send game join ACK to player: %v", err)
	}
}
