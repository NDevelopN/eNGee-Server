package utils

type Player struct {
	PID    string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Game struct {
	GID       string   `json:"id"`
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Status    string   `json:"status"`
	OldStatus string   `json:"old_status"`
	Leader    string   `json:"leader"`
	Rules     Rules    `json:"rules"`
	Players   []Player `json:"players"`
}

type Rules struct {
	Rounds     int    `json:"rounds"`
	MinPlrs    int    `json:"min_plrs"`
	MaxPlrs    int    `json:"max_plrs"`
	Timeout    int    `json:"timeout"`
	Additional string `json:"additional"`
}

type Join struct {
	PID string `json:"pid"`
	GID string `json:"gid"`
}

type GView struct {
	GID     string `json:"gid"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	CurPlrs int    `json:"cur_plrs"`
	MaxPlrs int    `json:"max_plrs"`
}

type GameMsg struct {
	Type    string `json:"type"`
	PID     string `json:"pid"`
	GID     string `json:"gid"`
	Content string `json:"content"`
}

type GameInfo struct {
	Games []GView `json:"games"`
}

type PlrList struct {
	Players []Player `json:"players"`
}

type SList struct {
	List []string `json:"list"`
}

type Pair struct {
	First  string `json:"first"`
	Second string `json:"second"`
}

type PairList struct {
	List []Pair `json:"list"`
}
type ACK struct {
	Message string `json:"message"`
}
