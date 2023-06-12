package consequences

import (
	"encoding/json"
	"log"
	"strings"

	db "Engee-Server/database"
	u "Engee-Server/utils"
)

func Start(gid string, broadcast func(string, []byte)) {
	p, err := json.Marshal(gMap[gid].Prompts)
	if err != nil {
		log.Printf("[Error] Failed to marshal prompt list: %v", err)
		return
	}

	db.UpdateGamePlayerStatus(gid, "Writing")

	msg, err := json.Marshal(u.GameMsg{
		Type:    "Prompts",
		PID:     "",
		GID:     gid,
		Content: string(p),
	})

	if err != nil {
		log.Printf("[Error] Failed to marshal prompt message: %v", err)
		return
	}

	broadcast(gid, msg)
}

var HandleInput u.GHandler = func(msg u.GameMsg, broadcast func(string, []byte)) {

	gm, err := db.GetGame(msg.GID)
	if err != nil {
		log.Printf("[Error] Failed to get game from db in gHandler: %v", err)
		return
	}

	if strings.ToLower(gm.Type) != "consequences" {
		log.Printf("[Error] Gametype mismatch: %v", gm.Type)
		return
	}

	switch msg.Type {
	case "Create":
		CreateGame(msg.GID, msg.Content)
	case "Connect":
		PlayerJoin(msg.GID, msg.PID)
	case "Leave":
		PlayerLeave(msg.GID, msg.PID)
	case "Start":
		Start(msg.GID, broadcast)
	case "Update":
		UpdateGame(msg.GID, msg.PID, msg.Content)
	case "Reply":
		HandleReply(msg)
	case "Restart":
		ResetGame(msg.GID)
	case "End":
		EndGame(msg.GID)
	default:
		log.Printf("[Error] No matching message type: %v", msg.Type)
	}
}

// TODO simplify reply structs
func HandleReply(msg u.GameMsg) {
	var r Replies

	err := json.Unmarshal([]byte(msg.Content), &r)
	if err != nil {
		log.Printf("[Error] Failed to unmarshal replies: %v", err)
	}

	length := len(gMap[msg.GID].Prompts)

	if len(r.List) != length {
		log.Printf("[Error] Mismatch in length of replies and prompts: %v : %v", len(r.List), length)
		return
	}

	AddStory(msg.GID, msg.PID, r.List)

	plr, err := db.GetPlayer(msg.PID)
	if err != nil {
		log.Printf("[Error] Could not get player from db for congame: %v", err)
		return
	}

	plr.Status = "Submitted"

	err = db.UpdatePlayer(plr)
	if err != nil {
		log.Printf("[Error] Could not update player for congame: %v", err)
		return
	}

	plrs, err := db.GetGamePlayers(msg.GID)
	if err != nil {
		log.Printf("[Error] Could not get game players for accept message: %v", err)
		return
	}

	plrList, err := json.Marshal(plrs)
	if err != nil {
		log.Printf("[Error] Could not marshal game players for accept message: %v", err)
		return
	}

	u.SockSend(u.Connections[msg.GID][msg.PID], "Accept", msg.GID, msg.PID, string(plrList))

	CheckComplete(msg.GID)
}

func AddStory(gid string, pid string, s []string) {
	gRef := gMap[gid]

	index := gRef.PMap[pid]
	gRef.Stories[index] = s
	gRef.ReadyCount++

	gMap[gid] = gRef
}

func CheckComplete(gid string) {
	gRef := gMap[gid]

	if gRef.ReadyCount >= gRef.PlayerCount {
		gRef.Stories = ShuffleStories(gRef.Stories)
		gMap[gid] = gRef
		EndRound(gid)
	}
}

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

func EndRound(gid string) {
	gRef := gMap[gid]

	for i, p := range gRef.PMap {
		var pl Story

		err := db.UpdateGamePlayerStatus(gid, "Reading")
		if err != nil {
			log.Printf("[Error] Could not update all players to reading status: %v", err)
		}

		for j := range gRef.Prompts {
			pl.Lines = append(pl.Lines, Line{Prompt: gRef.Prompts[j], Story: gRef.Stories[p][j]})
		}

		msg, err := json.Marshal(pl)
		if err != nil {
			log.Printf("[Error] Failed to marshal story: %v", err)
			return
		}

		u.SockSend(u.Connections[gid][i], "Story", gid, i, string(msg))
	}
}
