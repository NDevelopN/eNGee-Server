package lobby

import (
	"Engee-Server/room"
	"Engee-Server/user"
	"fmt"
)

var lobbies = make(map[string][]string)

func JoinUserToRoom(uid string, rid string) error {
	err := checkUserAndRoomExist(uid, rid)
	if err != nil {
		return err
	}

	if checkRoomLobbyExists(rid) {
		if checkRoomContainsUser(uid, rid) {
			return fmt.Errorf("user already in this room")
		}
	} else {
		lobbies[rid] = make([]string, 0)
	}

	lobbies[rid] = append(lobbies[rid], uid)

	return nil
}

func checkUserAndRoomExist(uid string, rid string) error {
	_, err := user.GetUser(uid)
	if err != nil {
		return err
	}

	_, err = room.GetRoom(rid)
	if err != nil {
		return err
	}

	return nil
}

func checkRoomLobbyExists(rid string) bool {
	_, found := lobbies[rid]
	return found
}

func checkRoomContainsUser(uid string, rid string) bool {
	for _, userID := range lobbies[rid] {
		if userID == uid {
			return true
		}
	}

	return false
}
