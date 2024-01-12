package user

import (
	"fmt"

	"Engee-Server/utils"

	"github.com/google/uuid"
)

type User struct {
	UID    string `json:"uid"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

var users = make(map[string]User)

func CreateUser(name string) (string, error) {
	err := utils.ValidateInputRefuseEmpty(name, nil)
	if err != nil {
		return "", err
	}

	var newUser User
	newUser.UID = uuid.NewString()
	newUser.Name = name
	newUser.Status = "New"

	users[newUser.UID] = newUser
	return newUser.UID, nil
}

func GetUser(uid string) (User, error) {
	return getUserByID(uid)
}

func UpdateUserName(uid string, name string) error {
	err := utils.ValidateInputRefuseEmpty(name, nil)
	if err != nil {
		return err
	}

	user, err := getUserByID(uid)
	if err != nil {
		return err
	}

	user.Name = name
	users[uid] = user

	return nil
}

func UpdateUserStatus(uid string, status string) error {
	err := utils.ValidateInputRefuseEmpty(status, nil)
	if err != nil {
		return err
	}

	user, err := getUserByID(uid)
	if err != nil {
		return err
	}

	user.Status = status
	users[uid] = user

	return nil
}

func DeleteUser(uid string) error {
	_, err := getUserByID(uid)
	if err != nil {
		return err
	}

	delete(users, uid)

	return nil
}

func getUserByID(uid string) (User, error) {
	var err error

	user, found := users[uid]
	if !found {
		err = fmt.Errorf("no user found with id: %q", uid)
	}

	return user, err
}
