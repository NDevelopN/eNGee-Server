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
