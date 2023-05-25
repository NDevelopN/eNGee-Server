package consequences

import (
	"encoding/json"
	"log"
	"strings"

	u "Engee-Server/utils"
)

func Start(gid string) {
	p, err := json.Marshal(Prompts{List: gMap[gid].Prompts})
	if err != nil {
		log.Printf("[Error] Failed to marshal prompt list: %v", err)
		return
	}

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

	u.Broadcast(gid, msg)
}

var HandleInput u.GHandler = func(msg u.GameMsg) {
	if strings.ToLower(u.Games[msg.GID].Type) != "consequences" {
		log.Printf("[Error] Gametype mismatch: %v", u.Games[msg.GID].Type)
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
		Start(msg.GID)
	case "Update":
		UpdateGame(msg.GID, msg.PID, msg.Content)
	case "Reply":
		HandleReply(msg)
	default:
		log.Printf("No matching message type: %v", msg.Type)
	}
}

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

	u.SockSend(u.Connections[msg.GID][msg.PID], "Accept", msg.GID, msg.PID, "")

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
		log.Printf("PID: %v", i)
		var pl Story

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
