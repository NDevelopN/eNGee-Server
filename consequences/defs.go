package consequences

import (
	"Engee-Server/utils"
	"encoding/json"
)

type ConSettings struct {
	Rounds  int      `json:"rounds"`
	Shuffle int      `json:"shuffle"`
	Timer1  int      `json:"timer1"`
	Timer2  int      `json:"timer2"`
	Prompts []string `json:"prompts"`
}

type ConVars struct {
	State    string
	SusState string
	Round    int
	Timer    int
	Settings ConSettings
	Stories  map[string][]string
}

var CVars map[string]ConVars

var defPrompts = []string{
	"Name of First Character",
	"Name of Second Character",
	"Location of the Scene",
	"Character 1's Action",
	"Character 2's Action",
	"Consequences of the Scene",
}

var defStory = []string{
	"Character 1",
	"Character 2",
	"Location",
	"Action 1",
	"Action 2",
	"Consequence",
}

var testSettings = ConSettings{
	Rounds:  1,
	Shuffle: 1,
	Timer1:  10,
	Timer2:  10,
	Prompts: defPrompts,
}

var defSettings = ConSettings{
	Rounds:  1,
	Shuffle: 1,
	Timer1:  0,
	Timer2:  0,
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
