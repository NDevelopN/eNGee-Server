package gameRegistry

import (
	"Engee-Server/utils"
	"fmt"
)

var urlRegistry = make(map[string]string)

func RegisterGameType(name string, url string) error {
	if name == "" {
		return fmt.Errorf("game type name is empty string")
	}

	err := utils.ValidateURL(url)
	if err != nil {
		return err
	}

	_, found := urlRegistry[name]
	if found {
		return fmt.Errorf("a gametype with that name already exists")
	}

	urlRegistry[name] = url
	return nil
}

func RemoveGameType(name string) error {
	_, found := urlRegistry[name]
	if !found {
		return fmt.Errorf("no matching gametype found")
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
		return "", fmt.Errorf("no matching gametype found")
	}

	return url, nil
}
