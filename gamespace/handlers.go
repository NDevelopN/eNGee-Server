package gamespace

import (
	consequences "Engee-Server/consequences"
	utils "Engee-Server/utils"
	"errors"
)

type HandlerFunc func(msg utils.GameMsg) (utils.GameMsg, error)

func TestHandler(msg utils.GameMsg) (utils.GameMsg, error) {
	switch msg.Type {
	case "Init":
		break
	case "Start":
		break
	case "Reset":
		break
	case "End":
		break
	case "Pause":
		break
	case "Remove":
		break
	case "Rules":
		break
	case "Status":
		break
	case "Leave":
		break
	default:
		return utils.ReplyError(msg, errors.New("invalid message Type"))
	}

	return utils.ReplyACK(msg, "TestMessage Accepted")
}

var TypeHandlers = map[string]HandlerFunc{
	"test":         TestHandler,
	"consequences": consequences.Handle,
}
