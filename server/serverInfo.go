package server

type Player struct {
	Name   string
	Games  map[string]string
	Status string
}

type Game struct {
	Name        string
	GameType    string
	MinPlayers  int
	MaxPlayers  int
	Status      string
	PlayerCount int
	Players     []string
}

var Plrs = map[string]Player{}

var Games = map[string]Game{}
