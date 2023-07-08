package consequences

import (
	"Engee-Server/utils"
	"encoding/json"
)

type ConSettings struct {
	Rounds  int `json:"rounds"`
	Shuffle int `json:"shuffle"`

	Prompts []string `json:"prompts"`
}

type Stories map[string][]string

type ConVars struct {
	settings ConSettings
	stories  Stories
}

type gVars map[string]ConVars

var defPrompts = []string{
	"Name of First Character",
	"Name of Second Character",
	"Location of the Scene",
	"Character 1's Action",
	"Character 2's Action",
	"Consequences of the Scene",
}

var testSettings = ConSettings{
	Rounds:  1,
	Shuffle: 1,
	Prompts: defPrompts,
}

var ts, _ = json.Marshal(testSettings)

var defGame = utils.Game{
	GID:             "",
	Name:            "Con Test",
	Type:            "Consequences",
	Status:          "Lobby",
	OldStatus:       "",
	Leader:          "",
	MinPlrs:         3,
	MaxPlrs:         10,
	CurPlrs:         0,
	AdditionalRules: string(ts),
}

var defPlr = utils.User{
	UID:    "",
	GID:    "",
	Name:   "Con Tester",
	Status: "Ready",
}
