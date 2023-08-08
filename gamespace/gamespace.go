package gamespace

import (
	g "Engee-Server/game"
	h "Engee-Server/handlers"
	u "Engee-Server/user"
	"Engee-Server/utils"
	"encoding/json"
	"time"

	"fmt"
	"log"
)

var Shutdown = map[string](chan int){}

func CleanUp(gid string) {
	Shutdown[gid] <- 0
	utils.RemoveConnectionPool(gid)
}

func pListUpdateBC(gid string, plrs []utils.User) error {
	pList, err := json.Marshal(plrs)
	if err != nil {
		return fmt.Errorf("could not marshal player list: %v", err)
	}

	msg := utils.GameMsg{
		GID:     gid,
		Type:    "Players",
		Content: string(pList),
	}

	err = utils.Broadcast(msg)
	if err != nil {
		return fmt.Errorf("could not broadcast player list: %v", err)
	}

	return nil
}

func activePlayersStatusUpdate(plrs []utils.User, status string) error {
	for _, plr := range plrs {
		if plr.Status == "Leaving" || plr.Status == "Spectating" {
			continue
		}

		plr.Status = status

		err := u.UpdateUser(plr)
		if err != nil {
			return fmt.Errorf("failed to update user [%s]: %v", plr.UID, err)
		}
	}

	return nil
}

func gameUpdateBC(game utils.Game) error {
	gm, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("could not marshal game update: %v", err)
	}

	upd := utils.GameMsg{
		GID:     game.GID,
		Type:    "Update",
		Content: string(gm),
	}

	err = utils.Broadcast(upd)
	if err != nil {
		return fmt.Errorf("could not broadcast game update: %v", err)
	}

	return nil
}

func gameStatusBC(gid string, status string) error {
	upd := utils.GameMsg{
		GID:     gid,
		Type:    "Status",
		Content: status,
	}

	err := utils.Broadcast(upd)
	if err != nil {
		return fmt.Errorf("could not broadcast game status update: %v", err)
	}

	return nil

}

func getValidPlrGame(msg utils.GameMsg) (utils.User, utils.Game, error) {
	var plr utils.User
	var game utils.Game

	plr, err := u.GetUser(msg.UID)
	if err != nil {
		return plr, game, fmt.Errorf("could not get user: %v", err)
	}

	if plr.GID != msg.GID {
		return plr, game, fmt.Errorf("player GID and msg GID do not match: %v", err)
	}

	game, err = g.GetGame(msg.GID)
	if err != nil {
		return plr, game, fmt.Errorf("could not get game: %v", err)
	}

	return plr, game, nil
}

func initialize(msg utils.GameMsg, game utils.Game) (string, string) {
	errStr := "[Error] Cannot initialize game: "

	log.Printf("Initializing game: %v", msg.GID)

	plrs, err := g.GetGamePlayers(msg.GID)
	if err != nil {
		end(msg, game)

		fmt.Printf("%v could not get game players: %v", errStr, err)
		return "Error", "No game players found."
	}

	pool, err := utils.GetConnections(msg.GID)
	if len(pool) == 0 || err != nil {
		time.Sleep(time.Second * 2)
		pool, err = utils.GetConnections(msg.GID)
		if len(pool) == 0 || err != nil {
			end(msg, game)
			fmt.Printf("%v no connections to game: %v", errStr, err)
			return "Error", "No available connections."
		}
	}

	handler, err := h.GetHandler(game.Type)
	if err != nil {
		log.Printf("%v could not get game handler: %v.", errStr, err)
		return "Error", "No game handler found for game type: " + game.Type + "."
	}

	cause, resp := handler(msg)
	if cause != "" {
		log.Printf("%v %v in game handler: %v.", errStr, cause, resp)
		return cause, resp
	}

	_, k := Shutdown[msg.GID]
	if k {
		Shutdown[msg.GID] <- 0
	}

	Shutdown[msg.GID] = make(chan int)

	go CheckGame(game)
	go CheckPlayers(game)

	game.Status = "Lobby"
	err = g.UpdateGame(game)
	if err != nil {
		fmt.Printf("%v could not update game to Lobby status: %v", errStr, err)
		return "Error", "Could not set game status."
	}

	err = activePlayersStatusUpdate(plrs, "Not Ready")
	if err != nil {
		fmt.Printf("%v could not update active players status: %v", errStr, err)
		return "Error", "Could not set active players' status."
	}

	plrs, err = g.GetGamePlayers(msg.GID)
	if err != nil {
		log.Printf("%v could not get game players after setting status: %v", errStr, err)
		return "Error", "Could not get active players' status."
	}

	pListUpdateBC(msg.GID, plrs)

	return "", ""
}

func GamespaceHandle(msg utils.GameMsg) (utils.GameMsg, error) {
	plr, game, err := getValidPlrGame(msg)
	if err != nil {
		log.Printf("[Error] Could not validate game or player: %v", err)
		return utils.CreateReply(msg, "Error", "Invalid ID(s) provided")
	}

	leader := game.Leader == plr.UID

	var cause, resp string

	switch msg.Type {
	case "Init":
		if leader {
			cause, resp = initialize(msg, game)
		} else {
			cause = "Error"
			resp = "Must be a leader to Init"
		}
	case "Start":
		if leader {
			cause, resp = start(msg, game)
		} else {
			cause = "Error"
			resp = "Must be a leader to Start"
		}
	case "Reset":
		if leader {
			cause, resp = reset(msg, game)
		} else {
			cause = "Error"
			resp = "Must be a leader to Reset"
		}
	case "End":
		if leader {
			cause, resp = end(msg, game)
		} else {
			cause = "Error"
			resp = "Must be a leader to End"
		}
	case "Pause":
		if leader {
			cause, resp = pause(msg, game)
		} else {
			cause = "Error"
			resp = "Must be a leader to Pause"
		}
	case "Remove":
		if leader {
			cause, resp = remove(msg, game)
		} else {
			cause = "Error"
			resp = "Must be a leader to Remove"
		}
	case "Status":
		cause, resp = status(msg, plr, game)
	case "Leave":
		log.Printf("Received leave")
		cause, resp = leave(msg, plr, game)
	default:
		errStr := "[Error] Cannot process " + msg.Type + " Request: "

		handler, err := h.GetHandler(game.Type)
		if err != nil {
			log.Printf("%v could not get game handler: %v.", errStr, err)
			cause = "Error"
			resp = "No game handler found for game type: " + game.Type
		} else {
			cause, resp = handler(msg)
			if cause != "" {
				log.Printf("%v %v in game handler: %v.", errStr, cause, resp)
			}
		}
	}

	if cause != "" {
		return utils.CreateReply(msg, cause, resp)
	}

	return utils.GameMsg{}, nil
}
