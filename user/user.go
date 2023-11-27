package user

import (
	"fmt"

	"github.com/google/uuid"
)

type user struct {
	UID    string
	Name   string
	Status string
}

var users = make(map[string]user)

func CreateUser(name string) (string, error) {
	err := validateUserName(name)
	if err != nil {
		return "", err
	}

	var newUser user
	newUser.UID = uuid.NewString()
	newUser.Name = name
	newUser.Status = "New"

	users[newUser.UID] = newUser
	return newUser.UID, nil
}

func validateUserName(name string) error {
	if name == "" {
		return fmt.Errorf("username is empty")
	}

	return nil
}

func GetUser(uid string) (user, error) {
	return getUserByID(uid)
}

func UpdateUserName(uid string, name string) error {
	err := validateUserName(name)
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

func getUserByID(uid string) (user, error) {
	var err error

	user, found := users[uid]
	if !found {
		err = fmt.Errorf("no user found with id: %q", uid)
	}

	return user, err
}
