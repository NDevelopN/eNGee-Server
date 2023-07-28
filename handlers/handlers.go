package handlers

import (
	db "Engee-Server/database"
	"Engee-Server/handlers/consequences"
	"Engee-Server/utils"
	"errors"
	"log"
)

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

	return utils.ReplyACK(msg), nil //, "TestMessage Accepted")
}

var typeHandlers = map[string]utils.HandlerFunc{
	"test":         TestHandler,
	"consequences": consequences.Handle,
}

func Init() {
	err := db.CreateGameTypes(typeHandlers)
	if err != nil {
		log.Fatalf("[Error] Failed to create game type list: %v", err)
	}
}

func GetHandlers() map[string]utils.HandlerFunc {
	return typeHandlers
}
