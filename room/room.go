package room

import (
	gameclient "Engee-Server/gameClient"
	"Engee-Server/utils"
	"encoding/json"
	"fmt"

	registry "Engee-Server/gameRegistry"

	"github.com/google/uuid"
)

type Room struct {
	RID    string `json:"rid"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Addr   string `json:"addr"`
}

var rooms = make(map[string]Room)

func CreateRoom(roomInfo []byte) (string, error) {
	var newRoom Room
	err := json.Unmarshal(roomInfo, &newRoom)
	if err != nil {
		return "", err
	}

	err = utils.ValidateInputRefuseEmpty(newRoom.Name, nil)
	if err != nil {
		return "", err
	}

	id := uuid.NewString()

	newRoom.RID = id
	newRoom.Status = "New"
	newRoom.Addr = ""

	rooms[id] = newRoom

	return id, nil
}

func GetRoom(rid string) (Room, error) {
	return getRoomByID(rid)
}

func GetRooms() (map[string]Room, error) {
	var err error = nil
	if len(rooms) == 0 {
		err = fmt.Errorf("no rooms to return")
	}

	return rooms, err
}

func GetRoomURL(rid string) (string, error) {
	room, err := getRoomByID(rid)
	if err != nil {
		return "", err
	}

	if room.Addr == "" {
		return "", fmt.Errorf("room URL not set")
	}

	return room.Addr, nil
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

	room.Addr, err = registry.GetGameURL(rType)
	if err == nil {
		rooms[rid] = room
		return nil
	} else {
		return err
	}
}

func CreateRoomInstance(rid string) error {
	room, err := getRoomByID(rid)
	if err != nil {
		return err
	}

	err = gameclient.CreateGame(rid, room.Addr)
	if err != nil {
		room.Status = "Created"
		rooms[rid] = room
	}

	return err
}

func DeleteRoom(rid string) error {
	_, err := getRoomByID(rid)
	if err != nil {
		return err
	}

	gameclient.EndGame(rid)

	delete(rooms, rid)

	return nil
}

func getRoomByID(rid string) (Room, error) {
	var err error

	room, found := rooms[rid]
	if !found {
		err = fmt.Errorf("no room found with id: %q", rid)
	}

	return room, err
}
