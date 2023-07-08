package gamespace

import (
	g "Engee-Server/game"
	u "Engee-Server/user"
	utils "Engee-Server/utils"
	"fmt"
)

func checkStatusCompatible(game utils.Game, status string) bool {

	switch game.Type {
	default:
		switch game.Status {
		//TODO what cases need to be handled?
		default:
			switch status {
			case "Ready":
				return true
			case "Not Ready":
				return true
			case "Joining":
				return true
			default:
				return false
			}
		}
	}
}

func ChangeStatus(pid string, gid string, status string) error {
	if status == "" {
		return fmt.Errorf("empty status provided")
	}

	plr, err := u.GetUser(pid)
	if err != nil {
		return fmt.Errorf("could not get matching player: %v", err)
	}

	game, err := g.GetGame(gid)
	if err != nil {
		return fmt.Errorf("could not get matching game: %v", err)
	}

	if !checkStatusCompatible(game, status) {
		return fmt.Errorf("current game status %v does not support submitted player status %v", game.Status, status)
	}

	plr.Status = status

	err = u.UpdateUser(plr)
	if err != nil {
		return fmt.Errorf("could not update user: %v", err)
	}

	return nil
}

func Leave(pid string, gid string) error {
	return nil
}
