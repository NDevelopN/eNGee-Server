package consequences

import (
	"Engee-Server/utils"
	"encoding/json"
)

const (
	LOBBY int = iota
	PROMPTS
	POSTPROMPTS
	STORIES
	POSTSTORIES
	ERROR
)

type ConSettings struct {
	Rounds  int      `json:"rounds"`
	Shuffle int      `json:"shuffle"`
	Timer1  int      `json:"timer1"`
	Timer2  int      `json:"timer2"`
	Prompts []string `json:"prompts"`
}

type ConVars struct {
	State    int
	Paused   bool
	Round    int
	Timer    int
	Settings ConSettings
	Active   int
	Ready    int
	Stories  map[string][]string
}

var CVars map[string]ConVars

var DefPrompts = []string{
	"Name of First Character",
	"Name of Second Character",
	"Location of the Scene",
	"Character 1's Action",
	"Character 2's Action",
	"Consequences of the Scene",
}

var DefStory = []string{
	"Character 1",
	"Character 2",
	"Location",
	"Action 1",
	"Action 2",
	"Consequence",
}

var TestSettings = ConSettings{
	Rounds:  1,
	Shuffle: 1,
	Timer1:  10,
	Timer2:  10,
	Prompts: DefPrompts,
}

var DefSettings = ConSettings{
	Rounds:  1,
	Shuffle: 1,
	Timer1:  0,
	Timer2:  0,
	Prompts: DefPrompts,
}

var Ts, _ = json.Marshal(TestSettings)

var DefGame = utils.Game{
	GID:             "",
	Name:            "Con Test",
	Type:            "consequences",
	Status:          "Lobby",
	OldStatus:       "",
	Leader:          "",
	MinPlrs:         3,
	MaxPlrs:         10,
	CurPlrs:         0,
	AdditionalRules: string(Ts),
}

var DefPlr = utils.User{
	UID:  "",
	GID:  "",
	Name: "Con Tester",
}
