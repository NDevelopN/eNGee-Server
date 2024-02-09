package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const configPath = "./config.json"

type Config struct {
	Port string `json:"server_port"`
}

func ReadConfig() Config {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("Could not read config file on launch: %w", err))
	}

	var payload Config
	err = json.Unmarshal(content, &payload)
	if err != nil {
		panic(fmt.Sprintf("Could not parse config file on launch: %w", err))
	}

	return payload
}
