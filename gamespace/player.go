package gamespace

import (
	c "Engee-Server/connections"
	h "Engee-Server/handlers"
	u "Engee-Server/user"
	"Engee-Server/utils"
	"encoding/json"
	"log"

	"golang.org/x/exp/slices"
)

var validStatus = []string{
	"Ready", "Not Ready", "Joining", "Leaving", "Spectating",
}

func join(msg utils.GameMsg, plr utils.User, game utils.Game) (string, string) {
	errStr := "[Error] Could not complete joining: "

	gm, err := json.Marshal(game)
	if err != nil {
		log.Printf("%v could not marshal game update: %v", errStr, err)
		return "Error", "Could not create game update reply."
	}

	rMsg := utils.GameMsg{
		UID:     plr.UID,
		GID:     game.GID,
		Type:    "Update",
		Content: string(gm),
	}

	err = c.SingleMessage(rMsg)
	if err != nil {
		log.Printf("%v could not send game update: %v", errStr, err)
		return "Error", "Could not send game update reply."
	}

	msg.Type = "Status"
	msg.Content = "Not Ready"

	cause, resp := status(msg, plr, game)
	if cause != "" {
		return cause, "Could not update status after joining game: " + resp
	}

	return "", ""
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

	if plr.Status != msg.Content {
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
	}

	return "", ""
}

func leave(msg utils.GameMsg, plr utils.User, game utils.Game) (string, string) {
	plr.GID = ""

	msg.Type = "Status"
	msg.Content = "Leaving"

	cause, resp := status(msg, plr, game)
	if cause != "" {
		return cause, "Could not remove player from game: " + resp
	}

	if game.CurPlrs == game.MinPlrs { //TODO add toggle for this
		msg.Type = "Reset"
		cause, resp = reset(msg, game)
		if cause != "" {
			log.Printf("[Error] Could not reset after falling below min players: %v", resp)
		}
	}

	return "", ""
}
