package gameRegistry

import (
	"Engee-Server/utils"
	"fmt"
)

var urlRegistry = make(map[string]string)

func RegisterGameMode(name string, url string) error {
	if name == "" {
		return fmt.Errorf("game mode name is empty string")
	}

	err := utils.ValidateURL(url)
	if err != nil {
		return err
	}

	_, found := urlRegistry[name]
	if found {
		return fmt.Errorf("a game mode with that name already exists")
	}

	urlRegistry[name] = url
	return nil
}

func RemoveGameType(name string) error {
	_, found := urlRegistry[name]
	if !found {
		return fmt.Errorf("no matching game mode found")
	}

	delete(urlRegistry, name)

	return nil
}

func GetGameTypes() []string {
	var gTypes []string
	for name := range urlRegistry {
		gTypes = append(gTypes, name)
	}

	return gTypes
}

func GetGameURL(name string) (string, error) {
	url, found := urlRegistry[name]
	if !found {
		return "", fmt.Errorf("no matching game mode found")
	}

	return url, nil
}
