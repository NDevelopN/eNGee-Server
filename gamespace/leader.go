package gamespace

import (
	g "Engee-Server/game"
	u "Engee-Server/user"
	utils "Engee-Server/utils"
	"fmt"
	"log"
)

func checkLeader(gid string, lid string) (utils.Game, error) {
	game, err := g.GetGame(gid)

	if err != nil {
		return game, fmt.Errorf("could not find game")
	}

	if game.Leader != lid {
		return game, fmt.Errorf("game leader and provided ID do not match")
	}

	return game, nil
}

func sendUpdate(gid string) error {
	//TODO

	return nil
}

func allPlayerStatusUpdate(plrs []utils.User, status string) error {
	for _, plr := range plrs {
		plr.Status = status

		err := u.UpdateUser(plr)
		if err != nil {
			return fmt.Errorf("failed to update user: %v", err)
		}
	}

	return nil
}

func Pause(gid string, lid string) error {
	game, err := checkLeader(gid, lid)
	if err != nil {
		return err
	}

	if game.Status == "Pause" {
		game.Status = game.OldStatus
		game.OldStatus = ""
	} else {
		game.OldStatus = game.Status
		game.Status = "Pause"
	}

	err = g.UpdateGame(game)
	if err != nil {
		return fmt.Errorf("failed to update game: %v", err)
	}

	err = sendUpdate(gid)
	if err != nil {
		return fmt.Errorf("failed to broadcast update: %v", err)
	}

	return nil
}

func Start(gid string, lid string) error {
	game, err := checkLeader(gid, lid)
	if err != nil {
		return err
	}

	if game.Status != "Lobby" && game.OldStatus != "Lobby" {
		return fmt.Errorf("cannot start game that is not in lobby state")
	}

	if game.CurPlrs < game.MinPlrs {
		return fmt.Errorf("game does not have enough players to start game")
	}

	plrs, err := g.GetGamePlayers(gid)
	if err != nil {
		return fmt.Errorf("could not get game players: %v", err)
	}

	//TODO is there a better way?
	ready := 0
	for i, plr := range plrs {
		log.Printf("player %d: %v", i, plr)
		if plr.Status == "Ready" {
			ready++
		}
	}

	//TODO add toggle for this requirement
	if ready <= game.CurPlrs/2 {
		return fmt.Errorf("less than half of the game players are ready (%d/%d)", ready, game.CurPlrs)
	}

	game.Status = "Play"
	game.OldStatus = ""

	err = allPlayerStatusUpdate(plrs, "Play")
	if err != nil {
		return err
	}

	err = g.UpdateGame(game)
	if err != nil {
		return fmt.Errorf("failed to update game: %v", err)
	}

	err = sendUpdate(gid)
	if err != nil {
		return fmt.Errorf("failed to broadcast update: %v", err)
	}

	return nil
}

func Reset(gid string, lid string) error {
	game, err := checkLeader(gid, lid)
	if err != nil {
		return err
	}

	game.Status = "Lobby"

	err = g.UpdateGame(game)
	if err != nil {
		return fmt.Errorf("failed to update game: %v", err)
	}

	plrs, err := g.GetGamePlayers(gid)
	if err != nil {
		return fmt.Errorf("could not get game players: %v", err)
	}
	err = allPlayerStatusUpdate(plrs, "Lobby")
	if err != nil {
		return err
	}

	err = sendUpdate(gid)
	if err != nil {
		return fmt.Errorf("failed to broadcast update: %v", err)
	}

	return nil
}

func End(gid string, lid string) error {
	game, err := checkLeader(gid, lid)
	if err != nil {
		return err
	}

	game.Status = "Ending"
	err = g.UpdateGame(game)
	if err != nil {
		return fmt.Errorf("failed to update game: %v", err)
	}

	err = sendUpdate(gid)
	if err != nil {
		return fmt.Errorf("failed to broadcast update: %v", err)
	}

	err = g.DeleteGame(gid)
	if err != nil {
		return fmt.Errorf("failed to delete game: %v", err)
	}

	return nil
}

func Rules(gid string, lid string, game utils.Game) error {
	return nil
}

func Remove(gid string, lid string, tid string) error {
	_, err := checkLeader(gid, lid)
	if err != nil {
		return err
	}

	if tid == lid {
		return fmt.Errorf("leader cannot remove themselves, must leave")
	}

	tUser, err := u.GetUser(tid)
	if err != nil {
		return fmt.Errorf("failed to get target user: %v", err)
	}

	tUser.GID = ""

	err = u.UpdateUser(tUser)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}
