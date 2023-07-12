package consequences

import (
	"Engee-Server/utils"
	"fmt"
)

func changeState(state string, gid string) {
	cVar := CVars[gid]
	cVar.State = state
	CVars[gid] = cVar
}

func initialize(msg utils.GameMsg) (utils.GameMsg, error) {
	return utils.GameMsg{}, nil
}

func start(msg utils.GameMsg) (utils.GameMsg, error) {
	return utils.GameMsg{}, nil
}

func reset(msg utils.GameMsg) (utils.GameMsg, error) {
	return utils.GameMsg{}, nil
}

func end(msg utils.GameMsg) (utils.GameMsg, error) {
	return utils.GameMsg{}, nil
}

func pause(msg utils.GameMsg) (utils.GameMsg, error) {
	return utils.GameMsg{}, nil
}

func rules(msg utils.GameMsg) (utils.GameMsg, error) {
	return utils.GameMsg{}, nil
}

func status(msg utils.GameMsg) (utils.GameMsg, error) {
	return utils.GameMsg{}, nil
}

func leave(msg utils.GameMsg) (utils.GameMsg, error) {
	return utils.GameMsg{}, nil
}

func reply(msg utils.GameMsg) (utils.GameMsg, error) {
	return utils.GameMsg{}, nil
}

func Handle(msg utils.GameMsg) (utils.GameMsg, error) {
	switch msg.Type {
	case "Init":
		return initialize(msg)
	case "Start":
		return start(msg)
	case "Reset":
		return reset(msg)
	case "End":
		return end(msg)
	case "Pause":
		return pause(msg)
	case "Remove":
		//Can be treated as leave from game mode point of view
		msg.UID = msg.Content
		return leave(msg)
	case "Rules":
		return rules(msg)
	case "Status":
		return status(msg)
	case "Leave":
		return leave(msg)
	case "Reply":
		return reply(msg)
	default:
		return utils.ReplyError(msg, fmt.Errorf("unknown message type: %v", msg.Type))
	}
}

func GetConState(gid string) (ConVars, error) {
	cv, f := CVars[gid]
	if !f {
		return cv, fmt.Errorf("could not find variables for consequences game: %v", gid)
	}

	return cv, nil
}
