package utils

type Player struct {
	Name   string
	Status string
	Games  map[string]string
}

type PlayerInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PlayerStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type PlayerList struct {
	Players []PlayerStatus `json:"players"`
}

type GameID struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Game struct {
	Name        string
	GameType    string
	Leader      string
	MinPlayers  int
	MaxPlayers  int
	Status      string
	PlayerCount int
	Players     map[string]string
}

// Limitted view of Games
type GameInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	GameType    string `json:"game_type"`
	Status      string `json:"status"`
	PlayerCount int    `json:"player_count"`
}

type GameList struct {
	Games []GameInfo `json:"games"`
}

// For creating / changing settings of a game
type GameRules struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	GameType   string `json:"gameType"`
	MinPlayers int    `json:"minPlayers"`
	MaxPlayers int    `json:"maxPlayers"`
	Additional string `json:"additional"`
}

// For Player/Game interaction (joining, leaving, ready, etc)
type GameOp struct {
	PID string `json:"pid"`
	GID string `json:"gid"`
}

type RemovePlr struct {
	AdminID string `json:"adminID"`
	GID     string `json:"gid"`
	PlrID   string `json:"plrID"`
}
