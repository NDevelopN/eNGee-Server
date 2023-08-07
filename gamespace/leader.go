package gamespace

import (
	g "Engee-Server/game"
	h "Engee-Server/handlers"
	u "Engee-Server/user"
	"Engee-Server/utils"
	"fmt"
	"log"
)

func pause(msg utils.GameMsg, game utils.Game) (string, string) {
	errStr := "[Error] Could not toggle pause"
	if game.Status == "Pause" {
		if game.OldStatus != "" {
			game.Status = game.OldStatus
			game.OldStatus = ""
		} else {
			log.Printf("%v paused game does not have old status.", errStr)
			return "Error", "Paused game does not have an old status to return to."
		}
	} else {
		game.OldStatus = game.Status
		game.Status = "Pause"
	}

	handler, err := h.GetHandler(game.Type)
	if err != nil {
		log.Printf("%v could not get game handler: %v.", errStr, err)
		return "Error", "No game handler found for game type: " + game.Type
	}

	cause, resp := handler(msg)
	if cause != "" {
		log.Printf("%v %v in game handler: %v.", errStr, cause, resp)
		return cause, resp
	}

	err = g.UpdateGame(game)
	if err != nil {
		log.Printf("%v could not update game: %v.", errStr, err)
		return "Error", "Could not apply update to game."
	}

	return "", ""
}

func start(msg utils.GameMsg, game utils.Game) (string, string) {
	warnStr := "[Warn] Cannot start game: "
	errStr := "[Error] Cannot start game: "

	if game.Status != "Lobby" {
		fmt.Printf("%v game status is not 'Lobby': %v.", warnStr, game.Status)
		return "Warn", "Cannot start a game that isn't in the Lobby."
	}

	if game.CurPlrs < game.MinPlrs {
		fmt.Printf("%v current player count too low: %d/%d", warnStr, game.CurPlrs, game.MinPlrs)
		return "Warn", "Cannot start a game without the minimum player count."
	}

	plrs, err := g.GetGamePlayers(game.GID)
	if err != nil {
		fmt.Printf("%v cannot find game players: %v", errStr, err)
		return "Error", "Could not find the game players."
	}

	ready := game.CurPlrs
	for _, plr := range plrs {
		if plr.Status == "Not Ready" || plr.Status == "Joining" {
			ready--
		}
	}

	if ready <= game.CurPlrs/2 { //TODO: Add toggle this
		fmt.Printf("%v not enough ready players: %v/%v", warnStr, ready, game.CurPlrs)
		return "Warn", "Cannot start a game with less than half of players ready."
	}

	game.Status = "Play"
	game.OldStatus = ""

	handler, err := h.GetHandler(game.Type)
	if err != nil {
		log.Printf("%v could not get game handler: %v.", errStr, err)
		return "Error", "No game handler found for game type: " + game.Type
	}

	cause, resp := handler(msg)
	if cause != "" {
		log.Printf("%v %v in game handler: %v", errStr, cause, resp)
		return cause, resp
	}

	err = activePlayersStatusUpdate(plrs, "Not Ready")
	if err != nil {
		fmt.Printf("%v could not update the status of the active players: %v.", errStr, err)
		return "Error", "Could not update players to 'Play' status."
	}

	err = g.UpdateGame(game)
	if err != nil {
		fmt.Printf("%v could not update the game: %v.", errStr, err)
		return "Error", "Could not apply game update."
	}

	return "", ""
}

func reset(msg utils.GameMsg, game utils.Game) (string, string) {
	errStr := "[Error] Cannot reset game: "

	game.Status = "Resetting"

	plrs, err := g.GetGamePlayers(game.GID)
	if err != nil {
		fmt.Printf("%v cannot find game players: %v.", errStr, err)
		return "Error", "Could not find the game players."
	}

	handler, err := h.GetHandler(game.Type)
	if err != nil {
		log.Printf("%v could not get game handler: %v.", errStr, err)
		return "Error", "No game handler found for game type: " + game.Type
	}

	cause, resp := handler(msg)
	if cause != "" {
		log.Printf("%v %v in game handler: %v", errStr, cause, resp)
		return cause, resp
	}

	err = activePlayersStatusUpdate(plrs, "Not Ready")
	if err != nil {
		fmt.Printf("%v could not update the status of the active players: %v.", errStr, err)
		return "Error", "Could not update players to 'Play' status."
	}

	err = g.UpdateGame(game)
	if err != nil {
		fmt.Printf("%v could not update the game: %v.", errStr, err)
		return "Error", "Could not apply game update."
	}

	msg.Type = "Init"

	return initialize(msg, game)
}

func end(msg utils.GameMsg, game utils.Game) (string, string) {
	errStr := "[Error] Cannot end game: "

	handler, err := h.GetHandler(game.Type)
	if err != nil {
		log.Printf("%v could not get game handler: %v.", errStr, err)
		return "Error", "No game handler found for game type: " + game.Type + "."
	}

	cause, resp := handler(msg)
	if cause != "" {
		log.Printf("%v %v in game handler: %v.", errStr, cause, resp)
		return cause, resp
	}

	eMsg := utils.GameMsg{
		GID:  msg.GID,
		Type: "End",
	}

	err = utils.Broadcast(eMsg)
	if err != nil {
		log.Printf("%v could not broadcast end message: %v", errStr, err)
		return "Error", "Could not broadcast game end message."
	}

	Shutdown[game.GID] <- 0

	err = g.DeleteGame(game.GID)
	if err != nil {
		fmt.Printf("%v could not delete the game: %v.", errStr, err)
		return "Error", "Could not finalize game deletion."
	}

	return "", ""
}

func remove(msg utils.GameMsg, game utils.Game) (string, string) {
	errStr := "[Error] Cannot remove player from game: "
	warnStr := "[Warn] Cannot remove player from game: "

	t := msg.Content
	if t == game.Leader {
		fmt.Printf("%v leader cannot remove themselves.", warnStr)
		return "Warn", "A leader cannot remove themselves, they must leave normally."
	}

	tUser, err := u.GetUser(t)
	if err != nil {
		fmt.Printf("%v failed to find target user: %v.", errStr, err)
		return "Error", "Could not find target user to remove."
	}

	tUser.GID = ""

	rMsg := utils.GameMsg{
		UID:  t,
		GID:  game.GID,
		Type: "Leave",
	}

	cause, resp := leave(rMsg, tUser, game)
	if cause != "" {
		log.Printf("%v %v in leave(removal): %v.", errStr, cause, resp)
		return cause, resp
	}

	rMsg.Type = "Removal"
	//TODO provide a reason?

	err = utils.SingleMessage(rMsg)
	if err != nil {
		log.Printf("%v could not update removed user: %v.", errStr, err)
		return "Error", "Could not send user removal notice."
	}

	return "", ""
}
