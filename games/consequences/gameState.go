package consequences

import (
	u "Engee-Server/utils"
	"encoding/json"
	"log"
)

var defaultPrompts = []string{
	"Character 1",
	"Character 2",
	"Location",
	"Character 1's action",
	"Character 2's action",
	"Consequences",
}

var gMap map[string]ConGame = make(map[string]ConGame)

func CreateGame(gid string, add string) {

	var p Prompts
	st := "default"

	//TODO maybe there's more to do here
	if add != "" {
		var s Settings
		err := json.Unmarshal([]byte(add), &s)
		if err != nil {
			log.Printf("[Error] Failed to unmarshal additional settings: %v", err)
			return
		}

		p.List = s.Prompts
		st = s.ShuffleType
	} else {
		p.List = defaultPrompts
	}

	gMap[gid] = ConGame{
		PlayerCount: 0,
		ReadyCount:  0,
		Prompts:     p.List,
		Stories:     make([][]string, 0),
		PMap:        make(map[string]int),
		Shuffle:     st,
	}
}

func ResetGame(gid string) {
	gRef := gMap[gid]
	gRef.Stories = make([][]string, len(gRef.Stories))
	gRef.ReadyCount = 0
	gMap[gid] = gRef

}

func EndGame(gid string) {
	delete(gMap, gid)
}

func PlayerJoin(gid string, pid string) {
	gRef := gMap[gid]
	_, k := gRef.PMap[pid]
	if k {
		log.Printf("[Error] Player already in the game")
		u.SockSend(u.Connections[gid][pid], "Error", gid, pid, "Player already in the game")
		return
	}

	gRef.PMap[pid] = len(gRef.Stories)
	gRef.Stories = append(gRef.Stories, make([]string, len(gRef.Prompts)))
	gRef.PlayerCount++

	gMap[gid] = gRef
}

func PlayerLeave(gid string, pid string) {
	gRef := gMap[gid]
	_, k := gRef.PMap[pid]
	if !k {
		log.Printf("[Error] Player not in the game")
		u.SockSend(u.Connections[gid][pid], "Error", gid, pid, "Player not in the game")
		return
	}

	//This should squash the slice
	key := gRef.PMap[pid]
	for ; key < len(gRef.Stories)-1; key++ {
		gRef.Stories[key] = gRef.Stories[key+1]
	}
	gRef.Stories[key] = nil

	gRef.PMap[pid] = -1
	gRef.PlayerCount--

	gMap[gid] = gRef

}

func UpdateGame(gid string, pid string, content string) {
	gm := gMap[gid]

	var s Settings
	err := json.Unmarshal([]byte(content), &s)
	if err != nil {
		log.Printf("[Error] Failed to unmarshal additional settings: %v", err)
		u.SockSend(u.Connections[gid][pid], "Error", gid, pid, "Could not read rules update")
		return
	}

	gm.Prompts = s.Prompts
	gm.Shuffle = s.ShuffleType

	gMap[gid] = gm

	u.SockSend(u.Connections[gid][pid], "ACK", gid, pid, "")
}
