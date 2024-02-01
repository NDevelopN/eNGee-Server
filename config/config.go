package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const configPath = "./config.json"

type Address struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type Config struct {
	Server Address `json:"server"`
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
