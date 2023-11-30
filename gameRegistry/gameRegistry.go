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

	_, found := registry[name]
	if found {
		return fmt.Errorf("a gametype with that name already exists")
	}

	registry[name] = buildFunc
	return nil
}

func RemoveGame(name string) error {
	_, found := registry[name]
	if !found {
		return fmt.Errorf("no matching gametype found")
	}

	delete(registry, name)

	return nil
}
