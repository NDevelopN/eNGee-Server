package gamespace

import (
	h "Engee-Server/handlers"
	u "Engee-Server/user"
	"Engee-Server/utils"
	"log"

	"golang.org/x/exp/slices"
)

var validStatus = []string{
	"Ready", "Not Ready", "Joining", "Leaving", "Spectating",
}

func status(msg utils.GameMsg, plr utils.User, game utils.Game) (string, string) {
	errStr := "[Error] Could not change player status: "

	if msg.Content == "" {
		log.Printf("%v empty status provided", errStr)
		return "Error", "No new status provided"
	}

	if !slices.Contains(validStatus, msg.Content) {
		log.Printf("%v invalid status provided: %v", errStr, msg.Content)
		return "Error", msg.Content + " is not a valid status"
	}

	plr.Status = msg.Content

	handler, err := h.GetHandler(game.Type)
	if err != nil {
		log.Printf("%v could not get game handler: %v", errStr, err)
		return "Error", "No game handler found for game type: " + game.Type
	}

	cause, resp := handler(msg)
	if cause != "" {
		log.Printf("Issue in game handler: %v", resp)
		return cause, resp
	}

	err = u.UpdateUser(plr)
	if err != nil {
		log.Printf("%v could not update user: %v", errStr, err)
		return "Error", "Could not apply update to user"
	}

	return "", ""
}

func tryEnd(msg utils.GameMsg, game utils.Game, errStr string) (string, string) {
	msg.Type = "End"
	cause, _ := end(msg, game)
	if cause != "" {
		Shutdown[game.GID] <- 0
		log.Printf("%v could not end game", errStr)
		return cause, "Game not ended"
	}

	return "", ""
}

func leave(msg utils.GameMsg, plr utils.User, game utils.Game) (string, string) {
	errStr := "[Error] Could not remove player from game: "

	plr.GID = ""

	msg.Type = "Status"
	msg.Content = "Leaving"

	cause, resp := status(msg, plr, game)
	if cause != "" {
		log.Printf("%v could not update status", errStr)
		return cause, "Could not remove player from game: " + resp

	}

	if game.CurPlrs == 1 {
		cause, resp = tryEnd(msg, game, errStr)
		if cause != "" {
			return cause, resp + " after removing last player."
		}
	}

	if game.CurPlrs == game.MinPlrs { //TODO add toggle for this
		msg.Type = "Reset"
		cause, _ = reset(msg, game)
		if cause != "" {
			log.Printf("%v could not reset game after CurPlrs fell below MinPlrs ", errStr)

			tryCause, tryResp := tryEnd(msg, game, errStr)
			if tryCause != "" {
				return tryCause, tryResp + " after failing minimum player reset"
			}

			return cause, "Game not reset after falling below minimum player count"
		}
	}

	return "", ""
}
