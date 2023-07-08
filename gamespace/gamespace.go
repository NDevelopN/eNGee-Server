package gamespace

import (
	utils "Engee-Server/utils"
	"fmt"

	g "Engee-Server/game"
	u "Engee-Server/user"
)

func replyError(msg utils.GameMsg, err error) (utils.GameMsg, error) {
	reply := utils.GameMsg{
		Type:    "Error",
		GID:     msg.GID,
		UID:     msg.UID,
		Content: "There was an issue wiht the " + msg.Type + " request",
	}

	return reply, fmt.Errorf("error handling %v request; %v", msg.Type, err)

}

type HandlerFunc func(msg utils.GameMsg, game utils.Game) (utils.GameMsg, error)

var typeHandlers = map[string]HandlerFunc{}

func GamespaceHandle(msg utils.GameMsg) (utils.GameMsg, error) {
	game, err := g.GetGame(msg.GID)
	if err != nil {
		return replyError(msg, err)
	}

	handler := typeHandlers[game.Type]
	if handler == nil {
		return replyError(msg, fmt.Errorf(`game type %q does not have a handler`, game.Type))
	}

	user, err := u.GetUser(msg.UID)
	if err != nil {
		return replyError(msg, err)
	}

	if user.GID != msg.GID {
		return replyError(msg, fmt.Errorf("user is not in game provided"))
	}

	switch msg.Type {
	default:
		return utils.GameMsg{}, nil
	}
}
