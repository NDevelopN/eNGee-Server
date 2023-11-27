package user

import (
	"fmt"

	"github.com/google/uuid"
)

func CreateUser(name string) (string, error) {
	err := validateUserName(name)
	if err != nil {
		return "", err
	}

	id := uuid.NewString()

	return id, nil
}

func validateUserName(name string) error {
	if name == "" {
		return fmt.Errorf("username is empty")
	}

	return nil
}
