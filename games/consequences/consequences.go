package consequences

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"

	l "Engee-Server/lobby"
	u "Engee-Server/utils"
)

var prompts []string
var defaultPrompts = []string{
	"Character 1",
	"Character 2",
	"Location",
	"Character 1's action",
	"Character 2's action",
	"Consequences",
}

type conGame struct {
	playerCount int
	readyCount  int
	stories     [][]string
	pMap        map[string]int
}

var gMap map[string]conGame = make(map[string]conGame)

var custom = false

// TODO create a database for custom prompts
func SetPrompts(p []string) {
	prompts = p
	custom = true
}

func GetPrompts() []string {
	if custom {
		return prompts
	}

	return defaultPrompts
}

func AddStory(g string, p string, s []string) {
	gRef := gMap[g]

	gRef.pMap[p] = len(gRef.stories)
	gRef.stories = append(gRef.stories, s)
	gRef.readyCount++

	if gRef.readyCount >= gRef.playerCount {
		gRef.stories = ShuffleStories(gRef.stories)
		gMap[g] = gRef
		endRound(g)
		return
	}

	gMap[g] = gRef
}

// TODO offer more random shuffle?
func ShuffleStories(stories [][]string) [][]string {
	var ns [][]string = make([][]string, len(stories))

	for plr := range stories {
		ns[plr] = make([]string, len(stories[plr]))
		for line := range stories[plr] {
			k := (plr + line + 1) % len(stories)
			ns[plr][line] = stories[k][line]
		}

	}

	return ns
}

func endRound(g string) {
	gRef := gMap[g]

	for i, p := range gRef.pMap {
		var pl u.PairList

		for j := range GetPrompts() {
			pl.List = append(pl.List, u.Pair{First: GetPrompts()[j], Second: gRef.stories[p][j]})
		}

		msg, err := json.Marshal(pl)
		if err != nil {
			log.Printf("Could not marshal story: %v", err)
			return
		}

		l.SingleWrite("Story", i, g, string(msg))
	}

}

var start l.StartFunc = func(gid string) {
	log.Print("Starting Consequences...	")

	p, err := json.Marshal(u.SList{List: GetPrompts()})
	if err != nil {
		log.Printf("Cannot marshal prompts: %v", err)
		return
	}

	msg, err := json.Marshal(u.GameMsg{
		Type:    "Prompts",
		PID:     "",
		GID:     gid,
		Content: string(p),
	})
	if err != nil {
		log.Printf("Could not marshal start message: %v", err)
		return
	}

	pc := len(u.Games[gid].Players)

	cg := conGame{
		playerCount: pc,
		readyCount:  0,
		stories:     make([][]string, 0),
		pMap:        make(map[string]int),
	}

	gMap[gid] = cg

	l.UpdatePlayers(gid, msg)
}

var gameCon l.ConFunc = func(conn *websocket.Conn, pid string, gid string) {
	//TODO is there anything needed in this function right now?
}

var handler l.GHandler = func(m u.GameMsg) {
	if strings.ToLower(u.Games[m.GID].Type) != "con" {
		log.Printf("Gametype mismatch: %v", u.Games[m.GID].Type)
		return
	}

	switch m.Type {
	case "Prompts": // Handles creation of new set of prompts
		var p []string
		err := json.Unmarshal([]byte(m.Content), &p)
		if err != nil {
			log.Printf("Could not unmarshal prompts: %v", err)
			return
		}

		SetPrompts(p)

		l.SingleWrite("ACK", m.PID, m.GID, "")
		return

	case "Reply":
		// Handles replies to prompts
		var replies u.SList
		err := json.Unmarshal([]byte(m.Content), &replies)
		if err != nil {
			log.Printf("Could not unmarshal replies: %v", err)
			return
		}

		if len(replies.List) != len(GetPrompts()) {
			log.Printf("Mismatch in replies and prompts: \n%v. \n%v.\n", replies, GetPrompts())
			return
		}

		AddStory(m.GID, m.PID, replies.List)
		return

	case "Shuffle":
		// Accepts settigns for variations in the shuffling
		//TODO
	case "Save":
		// Saves the shuffled story for later viewing
		//TODO

	}
}

func Lobby(w http.ResponseWriter, r *http.Request) {
	l.Lobby(w, r, gameCon, start, handler)
}
