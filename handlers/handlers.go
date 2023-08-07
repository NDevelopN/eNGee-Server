package handlers

import (
	db "Engee-Server/database"
	"Engee-Server/utils"
	"fmt"
	"log"
)

func TestHandler(msg utils.GameMsg) (string, string) {
	switch msg.Type {
	case "Init":
		return "", ""
	case "Start":
		return "", ""
	case "Reset":
		return "", ""
	case "End":
		return "", ""
	case "Pause":
		return "", ""
	case "Remove":
		return "", ""
	case "Status":
		return "", ""
	case "Leave":
		return "", ""
	default:
		return "Error", "Unsupported message type: " + msg.Type
	}
}

var typeHandlers = map[string]utils.GHandler{
	"test": TestHandler,
	// "consequences": consequences.Handle,
}

func Init() {
	err := db.CreateGameTypes(typeHandlers)
	if err != nil {
		log.Fatalf("[Error] Failed to create game type list: %v", err)
	}
}

func GetHandlers() map[string]utils.GHandler {
	return typeHandlers
}

func GetHandler(gType string) (utils.GHandler, error) {
	h, k := typeHandlers[gType]
	if !k {
		return nil, fmt.Errorf("no handler registered for: %v", gType)
	}

	return h, nil
}
