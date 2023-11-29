package room

import (
	registry "Engee-Server/gameRegistry"
	"Engee-Server/utils"
	"fmt"

	"github.com/google/uuid"
)

type room struct {
	RID    string
	Name   string
	Type   string
	Status string
	Addr   string
}

var rooms = make(map[string]room)

func CreateRoom(name string) (string, error) {
	err := utils.ValidateInputRefuseEmpty(name, nil)
	if err != nil {
		return "", err
	}

	id := uuid.NewString()
	newRoom := room{
		RID:    id,
		Name:   name,
		Type:   "None",
		Status: "New",
	}

	rooms[id] = newRoom

	return id, nil
}

func GetRoom(rid string) (room, error) {
	return getRoomByID(rid)
}

func GetRooms() (map[string]room, error) {
	var err error = nil
	if len(rooms) == 0 {
		err = fmt.Errorf("no rooms to return")
	}

	return rooms, err
}

func UpdateRoomName(rid string, name string) error {
	err := utils.ValidateInputRefuseEmpty(name, nil)
	if err != nil {
		return err
	}

	room, err := getRoomByID(rid)
	if err != nil {
		return err
	}

	room.Name = name
	rooms[rid] = room

	return nil
}

func UpdateRoomStatus(rid string, status string) error {
	err := utils.ValidateInputRefuseEmpty(status, nil)
	if err != nil {
		return err
	}

	room, err := getRoomByID(rid)
	if err != nil {
		return err
	}

	room.Status = status
	rooms[rid] = room

	return nil
}

func UpdateRoomType(rid string, rType string) error {
	err := utils.ValidateInputRefuseEmpty(rType, nil)
	if err != nil {
		return err
	}

	room, err := getRoomByID(rid)
	if err != nil {
		return err
	}

	room.Type = rType
	rooms[rid] = room

	return nil
}

func BuildRoomGame(rid string) (string, error) {
	room, err := getRoomByID(rid)
	if err != nil {
		return "", err
	}

	if room.Type == "None" {
		return "", fmt.Errorf("no gametype set")
	}

	if room.Addr != "" {
		return "", fmt.Errorf("game already built. it must be closed before building again")
	}

	addr, err := registry.BuildGame(room.Type)
	if err == nil {
		room.Addr = addr
		rooms[rid] = room
	}

	return addr, err
}

func DeleteRoom(rid string) error {
	_, err := getRoomByID(rid)
	if err != nil {
		return err
	}

	delete(rooms, rid)

	return nil
}

func getRoomByID(rid string) (room, error) {
	var err error

	room, found := rooms[rid]
	if !found {
		err = fmt.Errorf("no room found with id: %q", rid)
	}

	return room, err
}
