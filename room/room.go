package room

import (
	"Engee-Server/utils"
	"fmt"

	"github.com/google/uuid"
)

type room struct {
	RID     string
	Name    string
	Type    string
	Status  string
	CurPlrs int
}

var rooms = make(map[string]room)

func CreateRoom(name string) (string, error) {
	err := utils.ValidateInputRefuseEmpty(name, nil)
	if err != nil {
		return "", err
	}

	id := uuid.NewString()
	newRoom := room{
		RID:     id,
		Name:    name,
		Type:    "None",
		Status:  "New",
		CurPlrs: 0,
	}

	rooms[id] = newRoom

	return id, nil
}

func GetRoom(rid string) (room, error) {
	return getRoomByID(rid)
}

func getRoomByID(rid string) (room, error) {
	var err error

	room, found := rooms[rid]
	if !found {
		err = fmt.Errorf("no room found with id: %q", rid)
	}

	return room, err
}
