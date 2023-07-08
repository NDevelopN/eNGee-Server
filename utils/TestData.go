package utils

import "encoding/json"

var DefUser = User{
	UID:    "",
	GID:    "",
	Name:   "Test name",
	Status: "",
}

var Dub, _ = json.Marshal(DefUser)
var Dus = string(Dub)

var DefLeader = User{
	UID:    "",
	GID:    "",
	Name:   "Leader",
	Status: "",
}

var Dlb, _ = json.Marshal(DefLeader)
var Dls = string(Dlb)

var DefGame = Game{
	GID:             "",
	Name:            "Test Game",
	Type:            "test",
	Status:          "Lobby",
	OldStatus:       "",
	Leader:          "",
	MinPlrs:         3,
	MaxPlrs:         5,
	CurPlrs:         0,
	AdditionalRules: "",
}

var Dgb, _ = json.Marshal(DefGame)
var Dgs = string(Dgb)

type Settings struct {
	TestSettings string   `json:"testsettings"`
	Ops          []string `json:"ops"`
}

var TestSettings = Settings{
	TestSettings: "test",
}
