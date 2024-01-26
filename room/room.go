package room

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"

	gameclient "Engee-Server/gameClient"
	registry "Engee-Server/gameRegistry"
	sErr "Engee-Server/stockErrors"
)

type Room struct {
	RID      string `json:"rid"`
	Name     string `json:"name"`
	GameMode string `json:"gamemode"`
	Status   string `json:"status"`
	Addr     string `json:"addr"`
}

var rooms = make(map[string]Room)

func CreateRoom(roomInfo []byte) (string, error) {
	var newRoom Room
	err := json.Unmarshal(roomInfo, &newRoom)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal room info: %w", err)
	}

	if newRoom.Name == "" {
		return "", &sErr.EmptyValueError{
			Field: "Name",
		}
	}

	id := uuid.NewString()

	newRoom.RID = id

	newRoom.Addr, err = registry.GetGamemodeURL(newRoom.GameMode)
	if err != nil {
		return "", fmt.Errorf("could not get get gamemode info: %w", err)
	}

	err = gameclient.CreateGameInstance(id, newRoom.Addr)
	if err != nil {
		return "", fmt.Errorf("could not create game instance: %w", err)
	}

	newRoom.Status = "Created"

	rooms[id] = newRoom

	return id, nil
}

func GetRoom(rid string) (Room, error) {
	if rid == "" {
		return Room{}, &sErr.EmptyValueError{
			Field: "RID",
		}
	}

	room, found := rooms[rid]
	if !found {
		return room, &sErr.MatchNotFoundError[string]{
			Space: "Rooms",
			Field: "RID",
			Value: rid,
		}
	}

	return room, nil
}

func GetRooms() []Room {
	return maps.Values(rooms)
}

func GetRoomURL(rid string) (string, error) {
	room, err := GetRoom(rid)
	if err != nil {
		return "", err
	}

	if room.Addr == "" {
		return "", &sErr.EmptyValueError{
			Field: "Addr",
		}
	}

	return room.Addr, nil
}

func UpdateRoomName(rid string, name string) error {
	if name == "" {
		return &sErr.EmptyValueError{
			Field: "Name",
		}
	}

	room, err := GetRoom(rid)
	if err != nil {
		return err
	}

	room.Name = name
	rooms[rid] = room

	return nil
}

func UpdateRoomStatus(rid string, status string) error {
	if status == "" {
		return &sErr.EmptyValueError{
			Field: "Status",
		}
	}

	room, err := GetRoom(rid)
	if err != nil {
		return err
	}

	room.Status = status
	rooms[rid] = room

	return nil
}

func UpdateRoomGameMode(rid string, roomGameMode string) error {
	if roomGameMode == "" {
		return &sErr.EmptyValueError{
			Field: "Gamemode",
		}
	}

	room, err := GetRoom(rid)
	if err != nil {
		return err
	}

	room.GameMode = roomGameMode

	room.Addr, err = registry.GetGamemodeURL(roomGameMode)
	if err != nil {
		return fmt.Errorf("could not get gamemode url from registry: %w", err)
	}

	rooms[rid] = room
	return nil
}

func InitializeRoomGame(rid string) error {
	room, err := GetRoom(rid)
	if err != nil {
		return err
	}

	err = gameclient.CreateGameInstance(rid, room.Addr)
	if err != nil {
		return fmt.Errorf("could not creat game instance: %w", err)
	}

	room.Status = "Created"
	rooms[rid] = room

	return nil
}

func DeleteRoom(rid string) error {
	_, err := GetRoom(rid)
	if err != nil {
		return err
	}

	err = gameclient.EndGame(rid)
	if err != nil {
		return fmt.Errorf("could not end game: %w", err)
	}

	delete(rooms, rid)

	return nil
}
