package gameRegistry

import (
	"fmt"
	"time"

	sErr "Engee-Server/stockErrors"
	"Engee-Server/utils"
)

var urlRegistry = make(map[string]string)
var heartbeats map[string]time.Time

func RegisterGameMode(name string, url string) error {
	if name == "" {
		return &sErr.EmptyValueError{
			Field: "Gamemode name",
		}
	}

	err := utils.ValidateURL(url)
	if err != nil {
		return fmt.Errorf("URL is invalid: %w", err)
	}

	_, found := urlRegistry[name]
	if found {
		return &sErr.MatchFoundError[string]{
			Space: "Gamemodes",
			Field: "Name",
			Value: name,
		}
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
	_, err := GetGamemodeURL(name)
	if err != nil {
		return err
	}

	heartbeats[name] = time.Now()

	return nil
}

func RemoveGameMode(name string) error {
	_, err := GetGamemodeURL(name)
	if err != nil {
		return err
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

func GetGamemodeURL(name string) (string, error) {
	if name == "" {
		return "", &sErr.EmptyValueError{
			Field: "Name",
		}
	}

	url, found := urlRegistry[name]
	if !found {
		return "", &sErr.MatchNotFoundError[string]{
			Space: "Gamemodes",
			Field: "Name",
			Value: name,
		}
	}

	return url, nil
}
