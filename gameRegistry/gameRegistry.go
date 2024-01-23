package gameRegistry

import (
	"Engee-Server/utils"
	"fmt"
	"time"
)

var urlRegistry = make(map[string]string)
var heartbeats map[string]time.Time

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

	if heartbeats == nil {
		heartbeats = make(map[string]time.Time)
		go utils.MonitorHeartbeats(&heartbeats, RemoveGameMode)
	}

	heartbeats[name] = time.Now()

	return nil
}

func Heartbeat(name string) error {
	_, found := urlRegistry[name]
	if !found {
		return fmt.Errorf("no game mode '%s' found", name)
	}

	heartbeats[name] = time.Now()

	return nil
}

func RemoveGameMode(name string) error {
	_, found := urlRegistry[name]
	if !found {
		return fmt.Errorf("no matching game mode found")
	}

	delete(urlRegistry, name)
	delete(heartbeats, name)

	return nil
}

func GetGameModes() []string {
	var gameModes []string
	for name := range urlRegistry {
		gameModes = append(gameModes, name)
	}

	return gameModes
}

func GetGameURL(name string) (string, error) {
	url, found := urlRegistry[name]
	if !found {
		return "", fmt.Errorf("no matching game mode found")
	}

	return url, nil
}
