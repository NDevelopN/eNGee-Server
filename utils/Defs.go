package utils

import "errors"

var NO_HANDLER = false

// type GHandler func(msg GameMsg, broadcast func(string, []byte))
type GHandler func(msg GameMsg) (string, string)

type Message interface {
	[]User | User | []Game | Game | Join | GameMsg | Response | []string
}

var ErrNoGame = errors.New("no matching rows in game table")

type User struct {
	UID    string `json:"uid"`
	GID    string `json:"gid"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Game struct {
	GID             string `json:"gid"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	Status          string `json:"status"`
	OldStatus       string `json:"old_status"`
	Leader          string `json:"leader"`
	MinPlrs         int    `json:"min_plrs"`
	MaxPlrs         int    `json:"max_plrs"`
	CurPlrs         int    `json:"cur_plrs"`
	AdditionalRules string `json:"additional_rules"`
}

type Join struct {
	UID string `json:"uid"`
	GID string `json:"gid"`
}

type GameMsg struct {
	Type    string `json:"type"`
	UID     string `json:"uid"`
	GID     string `json:"gid"`
	Content string `json:"content"`
}

type Response struct {
	Cause   string `json:"cause"`
	Message string `json:"message"`
}

type ACK struct {
	Message string `json:"message"`
}
