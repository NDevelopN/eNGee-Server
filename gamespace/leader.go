package gamespace

import (
	g "Engee-Server/game"
	u "Engee-Server/user"
	utils "Engee-Server/utils"
	"encoding/json"
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

func allPlayerStatusUpdate(plrs []utils.User, status string) error {
	for _, plr := range plrs {
		plr.Status = status

		err := u.UpdateUser(plr)
		if err != nil {
			return fmt.Errorf("failed to update user: %v", err)
		}
	}

	pList, err := json.Marshal(plrs)
	if err != nil {
		return fmt.Errorf("failed to marshal player list for status update: %v", err)
	}

	upd := utils.GameMsg{
		Type:    "Players",
		GID:     plrs[0].GID,
		Content: string(pList),
	}

	err = utils.Broadcast(upd)
	if err != nil {
		return fmt.Errorf("could not broacast update: %v", err)
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

	upd := utils.GameMsg{
		Type:    "Status",
		GID:     gid,
		Content: game.Status,
	}

	err = utils.Broadcast(upd)
	if err != nil {
		return fmt.Errorf("failed to broadcast update: %v", err)
	}

	return nil
}

func Start(msg utils.GameMsg) (string, error) {
	game, err := checkLeader(msg.GID, msg.UID)
	if err != nil {
		return "", err
	}

	if game.Status != "Lobby" && game.OldStatus != "Lobby" {
		return "Cannot start a game that is not in the Lobby", utils.ErrWarn
	}

	if game.CurPlrs < game.MinPlrs {
		return "Cannot start a game that does not have the minimum number of players", utils.ErrWarn
	}

	plrs, err := g.GetGamePlayers(msg.GID)
	if err != nil {
		return "", fmt.Errorf("could not get game players: %v", err)
	}

	ready := 0
	for _, plr := range plrs {
		if plr.Status == "Ready" {
			ready++
		}
	}

	if ready <= game.CurPlrs/2 {
		return "Cannot Start a game with less than half of all players ready", utils.ErrWarn
	}

	game.Status = "Play"
	game.OldStatus = ""

	err = allPlayerStatusUpdate(plrs, "Play")
	if err != nil {
		return "", fmt.Errorf("failed to update game players' status: %v", err)
	}

	err = g.UpdateGame(game)
	if err != nil {
		return "", fmt.Errorf("failed to update game: %v", err)
	}

	upd := utils.GameMsg{
		Type:    "Status",
		GID:     msg.GID,
		Content: game.Status,
	}

	err = utils.Broadcast(upd)
	if err != nil {
		return "", fmt.Errorf("could not broadcast update: %v", err)
	}

	return "", nil
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
	err = allPlayerStatusUpdate(plrs, "Not Ready")
	if err != nil {
		return fmt.Errorf("failed to update game players' status: %v", err)
	}

	upd := utils.GameMsg{
		Type:    "Status",
		GID:     gid,
		Content: game.Status,
	}

	err = utils.Broadcast(upd)
	if err != nil {
		return fmt.Errorf("could not broadcast update: %v", err)
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

	upd := utils.GameMsg{
		Type: "End",
		GID:  gid,
	}

	err = utils.Broadcast(upd)
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
	_, err := checkLeader(gid, lid)
	if err != nil {
		return err
	}

	err = g.UpdateGame(game)
	if err != nil {
		return fmt.Errorf("failed to update game: %v", err)
	}

	gm, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("failed to marshal game update: %v", err)
	}

	upd := utils.GameMsg{
		Type:    "Update",
		GID:     gid,
		Content: string(gm),
	}

	err = utils.Broadcast(upd)
	if err != nil {
		return fmt.Errorf("failed to broadcast game update: %v", err)
	}

	return Reset(gid, lid)
}

func Remove(gid string, lid string, tid string) error {
	_, err := checkLeader(gid, lid)
	if err != nil {
		return err
	}

	if tid == lid {
		return fmt.Errorf("leaders cannot remove themselves")
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

	rMsg := utils.GameMsg{
		Type: "Removal",
		GID:  gid,
		UID:  tid,
	}

	err = utils.SingleMessage(rMsg)
	if err != nil {
		log.Printf("[Error] %v could not inform user of their removal: %v", UpdatePlayerList(gid), err)
	}

	return UpdatePlayerList(gid)
}
