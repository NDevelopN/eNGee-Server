package gameRegistry

import (
	"Engee-Server/utils"
	"fmt"
)

var registry = make(map[string]func() (string, error))

func RegisterGameType(name string, buildFunc func() (string, error)) error {
	err := utils.ValidateInputRefuseEmpty(name, nil)
	if err != nil {
		return err
	}

	if checkGameTypeRegistered(name) {
		return fmt.Errorf("a gametype with that name already exists")
	}

	registry[name] = buildFunc
	return nil
}

func checkGameTypeRegistered(name string) bool {
	_, found := registry[name]
	return found
}
