package lobby

import (
	"fmt"
	"log"

	"Engee-Server/room"
	sErr "Engee-Server/stockErrors"
	"Engee-Server/user"
	"Engee-Server/utils"
)

var lobbies = make(map[string][]string)

func JoinUserToRoom(uid string, rid string) error {
	err := checkUserAndRoomExist(uid, rid)
	if err != nil {
		return err
	}

	if checkRoomLobbyExists(rid) {
		if checkRoomContainsUser(uid, rid) {
			return &sErr.MatchFoundError[string]{
				Space: "Room Users",
				Field: "UID",
				Value: uid,
			}
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

	err = requireRoomLobby(rid)
	if err != nil {
		return err
	}

	if !checkRoomContainsUser(uid, rid) {
		return &sErr.MatchNotFoundError[string]{
			Space: "Room Users",
			Field: "UID",
			Value: uid,
		}
	}

	return removeUIDFromLobby(uid, rid)
}

func RemoveUserFromAllRooms(uid string) error {
	_, err := user.GetUser(uid)
	if err != nil {
		return fmt.Errorf("could not get user: %w", err)
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
		return fmt.Errorf("could not remove UID from slice: %w", err)
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
		return nil, fmt.Errorf("could not find room: %w", err)
	}

	err = requireRoomLobby(rid)
	if err != nil {
		return nil, err
	}

	var users []user.User
	for _, uid := range lobbies[rid] {
		user, err := user.GetUser(uid)
		if err != nil {
			log.Printf("[Error] Attempted to get user in lobby room list: %v", err)
			err = removeUIDFromLobby(uid, rid)
			if err != nil {
				return nil, err
			}

			continue
		}

		users = append(users, user)
	}

	if len(users) == 0 {

		err = &sErr.EmptySetError{
			Space: "Lobby",
			Field: "Users",
		}
	}

	return users, err
}

func GetRoomUserCount(rid string) (int, error) {
	_, err := room.GetRoom(rid)
	if err != nil {
		return 0, fmt.Errorf("could not get room: %w", err)
	}

	err = requireRoomLobby(rid)
	if err != nil {
		return 0, err
	}

	return len(lobbies[rid]), nil
}

func checkUserAndRoomExist(uid string, rid string) error {
	_, err := user.GetUser(uid)
	if err != nil {
		return fmt.Errorf("could not get user: %w", err)
	}

	_, err = room.GetRoom(rid)
	if err != nil {
		return fmt.Errorf("could not get room: %w", err)
	}

	return nil
}

func requireRoomLobby(rid string) error {
	if !checkRoomLobbyExists(rid) {
		return &sErr.MatchNotFoundError[string]{
			Space: "Lobbies",
			Field: "RID",
			Value: rid,
		}
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
