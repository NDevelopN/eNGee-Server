package gamespace

import (
	g "Engee-Server/game"
	u "Engee-Server/user"
	utils "Engee-Server/utils"
	"encoding/json"
	"fmt"
)

func ChangeStatus(pid string, gid string, status string) error {
	if status == "" {
		return fmt.Errorf("empty status provided")
	}

	plr, err := u.GetUser(pid)
	if err != nil {
		return fmt.Errorf("could not get matching player: %v", err)
	}

	plr.Status = status

	err = u.UpdateUser(plr)
	if err != nil {
		return fmt.Errorf("could not update user: %v", err)
	}

	return UpdatePlayerList(gid)
}

func Leave(pid string, gid string) error {
	plr, err := u.GetUser(pid)
	if err != nil {
		return fmt.Errorf("could not get matching player: %v", err)
	}

	if plr.GID != gid {
		return fmt.Errorf("mismatch between player GID [%v] and provided GID [%v]", plr.GID, gid)
	}

	_, err = g.GetGame(gid)
	if err != nil {
		return fmt.Errorf("could not get matching game: %v", err)
	}

	plr.GID = ""

	err = u.UpdateUser(plr)
	if err != nil {
		return fmt.Errorf("could not update player: %v", err)
	}

	plr.Status = "Leaving"

	pString, err := json.Marshal(plr)
	if err != nil {
		return fmt.Errorf("could not marshal player for message: %v", err)
	}

	msg := utils.GameMsg{
		Type:    "Player",
		GID:     gid,
		UID:     pid,
		Content: string(pString),
	}

	err = utils.SingleMessage(msg)
	if err != nil {
		return fmt.Errorf("could not send player removal message: %v", err)
	}

	err = utils.RemoveConnection(gid, plr.UID)
	if err != nil {
		return fmt.Errorf("could not remove connection: %v", err)
	}

	return UpdatePlayerList(gid)
}
