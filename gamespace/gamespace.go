package gamespace

import (
	utils "Engee-Server/utils"
	"fmt"

	g "Engee-Server/game"
	u "Engee-Server/user"
)

var errNotLeader = fmt.Errorf("player is not the game leader")

func replyError(msg utils.GameMsg, err error) (utils.GameMsg, error) {
	reply := utils.GameMsg{
		Type:    "Error",
		GID:     msg.GID,
		UID:     msg.UID,
		Content: "There was an issue wiht the " + msg.Type + " request",
	}

	return reply, fmt.Errorf("error handling %v request; %v", msg.Type, err)
}

func replyACK(msg utils.GameMsg) utils.GameMsg {
	return utils.GameMsg{
		Type: "ACK",
		GID:  msg.GID,
		UID:  msg.UID,
	}
}

type HandlerFunc func(msg utils.GameMsg, game utils.Game) (utils.GameMsg, error)

func testHandler(msg utils.GameMsg, game utils.Game) (utils.GameMsg, error) {
	return replyACK(msg), nil
}

var typeHandlers = map[string]HandlerFunc{
	"test": testHandler,
}

func initialize(msg utils.GameMsg, game utils.Game, handler HandlerFunc) (utils.GameMsg, error) {
	//TODO is there any generic Gamespace initalization?

	return handler(msg, game)
}

func start(msg utils.GameMsg, game utils.Game, handler HandlerFunc) (utils.GameMsg, error) {
	err := Start(msg.GID, msg.UID)
	if err != nil {
		return replyError(msg, err)
	}

	return handler(msg, game)
}

func reset(msg utils.GameMsg, game utils.Game, handler HandlerFunc) (utils.GameMsg, error) {
	err := Reset(msg.GID, msg.UID)
	if err != nil {
		return replyError(msg, err)
	}

	return handler(msg, game)
}

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

	leader := game.Leader == msg.UID

	switch msg.Type {
	case "Init":
		if leader {
			return initialize(msg, game, handler)
		} else {
			return replyError(msg, errNotLeader)
		}
	case "Start":
		if leader {
			return start(msg, game, handler)
		} else {
			return replyError(msg, errNotLeader)
		}
	case "Reset":
		if leader {
			return reset(msg, game, handler)
		} else {
			return replyError(msg, errNotLeader)
		}

	default:
		return replyError(msg, fmt.Errorf("unknown message type: %v", msg.Type))
	}
}
