package utils

import (
	"log"

	"github.com/gorilla/websocket"
)

var Plrs = map[string]Player{}

var Games = map[string]Game{}

var Connections = map[string](map[string]*websocket.Conn){}

func RemoveGame(gid string) {
	delete(Games, gid)
	delete(Connections, gid)
}

func RemovePlayer(gid string, pid string) string {
	gm := Games[gid]

	for i, p := range gm.Players {
		if p.PID == pid {
			gm.Players[i] = gm.Players[len(gm.Players)-1]
			gm.Players = gm.Players[:len(gm.Players)-1]
			break
		}
	}

	if len(gm.Players) <= 0 {
		RemoveGame(gid)
		return ""
	} else {
		if gm.Leader == pid {
			gm.Leader = gm.Players[0].PID
		}

		delete(Connections[gid], pid)
		Games[gid] = gm
		return gm.Leader
	}

}

func CheckForPlayer(pid string) (bool, Player) {

	var b bool = true

	p, k := Plrs[pid]
	if !k {
		log.Printf("[Error] Failed to find a matching player: %v", pid)
		b = false
	}

	return b, p
}

func CheckForGame(gid string) (bool, Game) {

	var b bool = true

	g, k := Games[gid]
	if !k {
		log.Printf("[Error] Failed to find a matching game: %v", gid)
		b = false
	}

	return b, g
}

func CheckGameForPlayer(gm Game, pid string) bool {
	for _, p := range gm.Players {
		if p.PID == pid {
			return true
		}
	}

	return false
}
