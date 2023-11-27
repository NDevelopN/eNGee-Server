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
	err := validateInputNoEmpty(name, nil)
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

func GetUser(uid string) (user, error) {
	return getUserByID(uid)
}

func UpdateUserName(uid string, name string) error {
	err := validateInputNoEmpty(name, nil)
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
	err := validateInputNoEmpty(status, nil)
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

func validateInputNoEmpty(input string, allowed map[string]struct{}) error {
	if input == "" {
		return fmt.Errorf("input is empty")
	}

	return validateInput(input, allowed)
}

func validateInput(input string, allowed map[string]struct{}) error {
	if len(allowed) == 0 {
		return nil
	}

	_, contains := allowed[input]
	if contains {
		return nil
	}

	return fmt.Errorf("%q is not a valid input", input)
}

func DeleteUser(uid string) error {
	_, err := getUserByID(uid)
	if err != nil {
		return err
	}

	delete(users, uid)

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
