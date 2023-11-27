package room

import (
	"Engee-Server/utils"

	"github.com/google/uuid"
)

type game struct {
	GID     string
	Name    string
	Type    string
	Status  string
	CurPlrs int
}

var games = make(map[string]game)

func CreateRoom(name string) (string, error) {
	err := utils.ValidateInputRefuseEmpty(name, nil)
	if err != nil {
		return "", err
	}

	id := uuid.NewString()
	newGame := game{
		GID:     id,
		Name:    name,
		Type:    "None",
		Status:  "New",
		CurPlrs: 0,
	}

	games[id] = newGame

	return id, nil
}
