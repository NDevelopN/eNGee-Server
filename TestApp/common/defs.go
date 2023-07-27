package common

import "Engee-Server/utils"

var Game = utils.Game{
	GID:             "",
	Name:            "TestGame",
	Type:            "test",
	Status:          "Lobby",
	OldStatus:       "",
	Leader:          "",
	MinPlrs:         1,
	MaxPlrs:         5,
	CurPlrs:         0,
	AdditionalRules: "",
}

var User = utils.User{
	UID:    "",
	GID:    "",
	Name:   "TestLeader",
	Status: "",
}
