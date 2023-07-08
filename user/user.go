package user

import (
	db "Engee-Server/database"
	utils "Engee-Server/utils"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func CreateUser(u utils.User) (string, error) {
	if u.Name == "" {
		return "", errors.New("user name is empty")
	}

	if u.GID != "" {
		return "", fmt.Errorf("a new user should not have a GID: %v", u.GID)
	}

	u.UID = uuid.NewString()

	err := db.CreateUser(u)
	if err != nil {
		return "", fmt.Errorf("could not create user in database: %v", err)
	}

	return u.UID, nil
}

func GetUser(uid string) (utils.User, error) {
	return utils.User{}, nil
}

func UpdateUser(n utils.User) error {
	return nil
}

func DeleteUser(uid string) error {
	return nil
}
