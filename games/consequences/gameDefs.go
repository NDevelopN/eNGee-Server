package consequences

type ConGame struct {
	PlayerCount int
	ReadyCount  int
	Prompts     []string
	Stories     [][]string
	PMap        map[string]int
	Shuffle     string
}

type ShuffleFunc func(stories [][]string) [][]string

type Prompts struct {
	List []string `json:"list"`
}

type Replies struct {
	List []string `json:"list"`
}

type Story struct {
	Lines []Line `json:"lines"`
}

type Line struct {
	Prompt string `json:"prompt"`
	Story  string `json:"story"`
}

type Settings struct {
	ShuffleType string   `json:"shuffle_type"`
	Prompts     []string `json:"prompts"`
}
