package gamespace

import (
	utils "Engee-Server/utils"
)

type HandlerFunc func(msg utils.GameMsg, game utils.Game) (utils.GameMsg, error)

func TestHandler(msg utils.GameMsg, game utils.Game) (utils.GameMsg, error) {
	return utils.ReplyACK(msg), nil
}

var TypeHandlers = map[string]HandlerFunc{
	"test": TestHandler,
}
