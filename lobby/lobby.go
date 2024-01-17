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

	if !checkRoomLobbyExists(rid) {
		return fmt.Errorf("lobby for room does not exist")
	}

	if !checkRoomContainsUser(uid, rid) {
		return fmt.Errorf("room does not contain user")
	}

	return removeUIDFromLobby(uid, rid)
}

func RemoveUserFromAllRooms(uid string) error {
	_, err := user.GetUser(uid)
	if err != nil {
		return err
	}

	for rid := range lobbies {
		if checkRoomContainsUser(uid, rid) {
			err = removeUIDFromLobby(uid, rid)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func removeUIDFromLobby(uid string, rid string) error {
	var err error = nil
	lobbies[rid], err = utils.RemoveElementFromSliceOrdered(lobbies[rid], uid)
	if err != nil {
		return err
	}

	if len(lobbies[rid]) == 0 {
		room.DeleteRoom(rid)
		delete(lobbies, rid)
	}

	return nil
}

func GetUsersInRoom(rid string) ([]user.User, error) {
	_, err := room.GetRoom(rid)
	if err != nil {
		return nil, err
	}

	if !checkRoomLobbyExists(rid) {
		return nil, fmt.Errorf("lobby for room does not exist")
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

	if !checkRoomLobbyExists(rid) {
		return 0, fmt.Errorf("lobby for room does not exist")
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
