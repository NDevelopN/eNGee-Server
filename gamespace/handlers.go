package gamespace

import (
	consequences "Engee-Server/consequences"
	utils "Engee-Server/utils"
)

type HandlerFunc func(msg utils.GameMsg) (utils.GameMsg, error)

func TestHandler(msg utils.GameMsg) (utils.GameMsg, error) {
	return utils.ReplyACK(msg), nil
}

var TypeHandlers = map[string]HandlerFunc{
	"test":         TestHandler,
	"consequences": consequences.Handle,
}
