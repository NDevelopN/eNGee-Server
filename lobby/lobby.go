package lobby

import (
	"Engee-Server/room"
	"Engee-Server/user"
	"Engee-Server/utils"
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

func RemoveUserFromRoom(uid string, rid string) error {
	err := checkUserAndRoomExist(uid, rid)
	if err != nil {
		return err
	}

	err = fmt.Errorf("room does not contain user")

	if checkRoomLobbyExists(rid) {
		if checkRoomContainsUser(uid, rid) {
			lobbies[rid], err = utils.RemoveElementFromSliceOrdered(lobbies[rid], uid)
		}
	}

	return err
}

func GetUsersInRoom(rid string) ([]user.User, error) {
	_, err := room.GetRoom(rid)
	if err != nil {
		return nil, err
	}

	var users []user.User
	for _, uid := range lobbies[rid] {
		user, err := user.GetUser(uid)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if len(users) == 0 {
		err = fmt.Errorf("no users in room")
	}

	return users, err
}

func GetRoomUserCount(rid string) (int, error) {
	_, err := room.GetRoom(rid)
	if err != nil {
		return 0, err
	}

	return len(lobbies[rid]), nil
}

func GetRoomLeader(rid string) (string, error) {
	_, err := room.GetRoom(rid)
	if err != nil {
		return "", err
	}

	if checkRoomLobbyExists(rid) {
		return lobbies[rid][0], nil
	}

	return "", fmt.Errorf("lobby for %q does not exists", rid)
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
