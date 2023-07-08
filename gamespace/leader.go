package gamespace

import (
	g "Engee-Server/game"
	utils "Engee-Server/utils"
	"fmt"
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
	return nil
}

func Reset(gid string, lid string) error {
	return nil
}

func End(gid string, lid string) error {
	return nil
}

func Rules(gid string, lid string, game utils.Game) error {
	return nil
}

func Remove(gid string, lid string, tid string) error {
	return nil
}
