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
