package user

import (
	db "Engee-Server/database"
	g "Engee-Server/game"
	utils "Engee-Server/utils"

	"errors"
	"fmt"

	"github.com/google/uuid"
)

func CreateUser(u utils.User) (string, error) {
	if u.Name == "" {
		return "", errors.New("user name is empty")
	}

	if u.UID != "" {
		return "", fmt.Errorf("a new user should not have a UID: %v", u.UID)
	}
	if u.GID != "" {
		return "", fmt.Errorf("a new user should not have a GID: %v", u.GID)
	}

	if u.Status != "" {
		return "", fmt.Errorf("a new user should not have a status: %v", u.Status)
	}

	u.UID = uuid.NewString()
	u.Status = "New"

	err := db.CreateUser(u)
	if err != nil {
		return "", fmt.Errorf("could not create user in database: %v", err)
	}

	return u.UID, nil
}

func GetUser(uid string) (utils.User, error) {
	return db.GetUser(uid)
}

func joinUserToGame(gid string, uid string) error {
	u, err := GetUser(uid)
	if err != nil {
		return fmt.Errorf("could not find user in database: %v", err)
	}

	if u.GID != "" {
		return fmt.Errorf("user already in a game: %v", u.GID)
	}

	game, err := g.GetGame(gid)
	if err != nil {
		return fmt.Errorf("could not find game in database: %v", err)
	}

	if game.CurPlrs == game.MaxPlrs {
		return fmt.Errorf("not enough space in game for new player: %v/%v", game.CurPlrs, game.MaxPlrs)
	}

	if game.Leader == "" {
		game.Leader = uid
		err = g.UpdateGame(game)
		if err != nil {
			return fmt.Errorf("could not update (empty) game leader: %v", err)
		}
	}

	err = g.ChangePlayerCount(game, 1)
	if err != nil {
		return fmt.Errorf("could not change game player count: %v", err)
	}

	u.GID = gid

	err = db.UpdateUser(u)
	if err != nil {
		return fmt.Errorf("could not update user in database: %v", err)
	}

	return nil
}

func removeGID(u utils.User) error {
	u.GID = ""
	err := db.UpdateUser(u)
	if err != nil {
		return fmt.Errorf("could not update user in database: %v", err)
	}

	return nil
}

func removeUserFromGame(gid string, uid string) error {
	u, err := GetUser(uid)
	if err != nil {
		return fmt.Errorf("could not find user in database: %v", err)
	}

	if u.GID != gid {
		return fmt.Errorf("user not in provided game: %v", u.GID)
	}

	game, err := g.GetGame(gid)
	if err != nil {
		nuErr := removeGID(u)
		return fmt.Errorf("%v - could not find matching game: %v", nuErr, err)
	}

	err = g.ChangePlayerCount(game, -1)
	if err != nil {
		return fmt.Errorf("could not change game palyer count :%v", err)
	}

	err = removeGID(u)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUser(n utils.User) error {
	o, err := GetUser(n.UID)
	if err != nil {
		return fmt.Errorf("could not get user to update: %v", err)
	}

	if n.GID != o.GID {
		if n.GID == "" {
			err = removeUserFromGame(o.GID, o.UID)
			if err != nil {
				return fmt.Errorf("could not remove user from game: %v", err)
			}
		} else {
			err = joinUserToGame(n.GID, n.UID)
			if err != nil {
				return fmt.Errorf("could not join user to game: %v", err)
			}
		}
	}

	if n.Name == "" {
		return fmt.Errorf("provided user name is empty")
	}

	if n.GID != "" {
		_, err := g.GetGame(n.GID)

		if err != nil {
			return fmt.Errorf("could not find a matching game for user GID: %v", err)
		}
	}

	err = db.UpdateUser(n)
	if err != nil {
		return fmt.Errorf("could not update user in database: %v", err)
	}

	return nil
}

func DeleteUser(uid string) error {
	u, err := GetUser(uid)
	if err != nil {
		return fmt.Errorf("could not get user from database: %v", err)
	}

	if u.GID != "" {
		err = removeUserFromGame(u.GID, uid)
		if err != nil {
			return fmt.Errorf("could not remove user from game: %v", err)
		}
	}

	err = db.RemoveUser(uid)
	if err != nil {
		return fmt.Errorf("could not delete user from database: %v", err)
	}
	return nil
}
