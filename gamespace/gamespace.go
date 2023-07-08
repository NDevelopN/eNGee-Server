package gamespace

import (
	utils "Engee-Server/utils"
	"encoding/json"
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

func end(msg utils.GameMsg, game utils.Game, handler HandlerFunc) (utils.GameMsg, error) {
	err := End(msg.GID, msg.UID)
	if err != nil {
		return replyError(msg, err)
	}

	return handler(msg, game)
}

func pause(msg utils.GameMsg, game utils.Game, handler HandlerFunc) (utils.GameMsg, error) {
	err := Pause(msg.GID, msg.UID)
	if err != nil {
		return replyError(msg, err)
	}

	return handler(msg, game)
}

func remove(msg utils.GameMsg, game utils.Game, handler HandlerFunc) (utils.GameMsg, error) {
	err := Remove(msg.GID, msg.UID, msg.Content)
	if err != nil {
		return replyError(msg, err)
	}

	return handler(msg, game)
}

func rules(msg utils.GameMsg, game utils.Game, handler HandlerFunc) (utils.GameMsg, error) {
	var gm utils.Game
	err := json.Unmarshal([]byte(msg.Content), &gm)
	if err != nil {
		return replyError(msg, err)
	}

	err = Rules(msg.GID, msg.UID, gm)
	if err != nil {
		return replyError(msg, err)
	}

	return handler(msg, game)
}

func status(msg utils.GameMsg, game utils.Game, handler HandlerFunc) (utils.GameMsg, error) {
	err := ChangeStatus(msg.UID, msg.GID, msg.Content)
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
	case "End":
		if leader {
			return end(msg, game, handler)
		} else {
			return replyError(msg, errNotLeader)
		}
	case "Pause":
		if leader {
			return pause(msg, game, handler)
		} else {
			return replyError(msg, errNotLeader)
		}
	case "Remove":
		if leader {
			return remove(msg, game, handler)
		} else {
			return replyError(msg, errNotLeader)
		}
	case "Rules":
		if leader {
			return rules(msg, game, handler)
		} else {
			return replyError(msg, errNotLeader)
		}
	case "Status":
		return status(msg, game, handler)
	default:
		return replyError(msg, fmt.Errorf("unknown message type: %v", msg.Type))
	}
}
