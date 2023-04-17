package server

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

type GameRules struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	GameType   string `json:"gameType"`
	MinPlayers int    `json:"minPlayers"`
	MaxPlayers int    `json:"maxPlayers"`
	Additional string `json:"additional"`
}

type Join struct {
	PID string `json:"pid"`
	GID string `json:"gid"`
}
