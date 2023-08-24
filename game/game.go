package game

import (
	c "Engee-Server/connections"
	db "Engee-Server/database"
	"Engee-Server/utils"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func checkType(gType string) bool {
	if utils.NO_HANDLER {
		return true
	}
	tList, err := db.GetGameTypes()
	if err != nil {
		log.Printf("[Error] could not get game types from database: %v", err)
		return false
	}

	for _, t := range tList {
		if t == gType {
			return true
		}
	}

	return false
}

func CreateGame(g utils.Game) (string, error) {
	if g.Name == "" {
		return "", fmt.Errorf("provided game name is empty")
	}

	if g.GID != "" {
		return "", fmt.Errorf("a new game must not have a GID: %v", g.GID)
	}

	if g.Type == "" {
		return "", fmt.Errorf("provided game type is empty")
	}

	if !checkType(g.Type) {
		return "", fmt.Errorf("provided game type (%v) is not supported", g.Type)
	}

	if g.MinPlrs > g.MaxPlrs {
		return "", fmt.Errorf("provided minPlrs is greater than provided maxPlrs")
	}

	if g.MinPlrs < 0 {
		return "", fmt.Errorf("provided minPlrs is negative")
	}

	if g.CurPlrs != 0 {
		return "", fmt.Errorf("a new game must have CurPlrs == 0: %v", g.CurPlrs)
	}

	if g.OldStatus != "" {
		return "", fmt.Errorf("a new game must not have an OldStatus: %v", g.OldStatus)
	}

	g.GID = uuid.NewString()
	g.Status = "Lobby"

	err := db.CreateGame(g)
	if err != nil {
		return "", fmt.Errorf("failed to create game in database: %v", err)
	}

	c.AddConnectionPool(g.GID)

	return g.GID, nil
}

func GetGames() ([]utils.Game, error) {
	return db.GetAllGames()
}

func GetGame(gid string) (utils.Game, error) {
	return db.GetGame(gid)
}

func GetGamePlayers(gid string) ([]utils.User, error) {
	_, err := db.GetGame(gid)
	if err == sql.ErrNoRows {
		return nil, utils.ErrNoGame
	} else if err != nil {
		return nil, fmt.Errorf("failed to find game in database: %v", err)
	}

	plrs, err := db.GetGamePlayers(gid)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("failed to get players from database: %v", err)
	}

	return plrs, nil
}

func UpdateGame(g utils.Game) error {
	og, err := db.GetGame(g.GID)
	if err != nil {
		return fmt.Errorf("cannot find game in database to update: %v", err)
	}

	if g.Name == "" {
		return fmt.Errorf("cannot set game name to empty string")
	}

	if g.Type == "" {
		return fmt.Errorf("cannot set game type to empty string")
	}

	if !checkType(g.Type) {
		return fmt.Errorf("provided game type (%v) is not supported", g.Type)
	}

	if g.Status == "" {
		return fmt.Errorf("cannot set game status to empty string")
	}

	if g.Leader == "" {
		return fmt.Errorf("cannot set game leader to empty string")
	}

	if og.CurPlrs != g.CurPlrs {
		return fmt.Errorf("cannot change curPlrs: Old (%v) New (%v)", og.CurPlrs, g.CurPlrs)
	}

	_, err = db.GetUser(g.Leader)
	if err != nil {
		return fmt.Errorf("cannot find user to match leader ID")
	}

	if g.MinPlrs < 0 {
		return fmt.Errorf("provided minPlrs is negative: %v", g.MinPlrs)
	}

	if g.MinPlrs > g.MaxPlrs {
		return fmt.Errorf("provided minPlrs is greater than provided maxPlrs")
	}

	if g.CurPlrs > g.MaxPlrs {
		return fmt.Errorf("provided maxPlrs is less than curPlrs")
	}

	err = db.UpdateGame(g)
	if err != nil {
		return fmt.Errorf("could not update game in database: %v", err)
	}

	return nil
}

func ChangePlayerCount(g utils.Game, d int) error {
	g.CurPlrs += d

	if g.CurPlrs > g.MaxPlrs {
		return fmt.Errorf("the game is too full: %v/%v", g.CurPlrs, g.MaxPlrs)
	}

	if g.CurPlrs < g.MinPlrs && g.Status != "Lobby" {
		g.Status = "Lobby"
	}

	if g.CurPlrs <= 0 {
		err := DeleteGame(g.GID)
		if err != nil {
			return fmt.Errorf("could not delete empty game: %v", err)
		}

		return nil
	}

	err := UpdateGame(g)
	if err != nil {
		return fmt.Errorf("could not update game: %v", err)
	}

	return nil
}

func DeleteGame(gid string) error {
	plrs, err := db.GetGamePlayers(gid)
	if err != nil {
		log.Printf("[Warn] Deleting game -- Could not get player from database: %v", err)
	}

	for _, p := range plrs {
		p.GID = ""
		err = db.UpdateUser(p)
		if err != nil {
			return fmt.Errorf("could not update player (clearing GID) in database: %v", err)
		}
	}

	err = db.RemoveGame(gid)
	if err != nil {
		return fmt.Errorf("could not delete the game from database: %v", err)
	}

	return nil
}

func JoinGame(gid string, uid string) error {
	game, err := db.GetGame(gid)
	if err != nil {
		return fmt.Errorf("could not find game in database: %v", err)
	}

	if game.CurPlrs >= game.MaxPlrs {
		return fmt.Errorf("not enough space in game for new player: %v/%v", game.CurPlrs, game.MaxPlrs)
	}

	if game.Leader == "" {
		game.Leader = uid
	}

	game.CurPlrs++
	err = db.UpdateGame(game)
	if err != nil {
		return fmt.Errorf("could not update game: %v", err)
	}

	return nil
}

func LeaveGame(gid string, uid string) error {
	game, err := db.GetGame(gid)

	if err != nil {
		return fmt.Errorf("could not find game in database: %v", err)
	}

	plrs, err := GetGamePlayers(gid)
	if err != nil || len(plrs) <= 1 {
		err = DeleteGame(gid)
		if err != nil {
			log.Printf("[Error] Could not delete game when all players left: %v", err)
		}
		return nil
	}

	if game.Leader == uid {
		leader := ""

		for _, p := range plrs {
			if p.UID != uid {
				leader = p.UID
				break
			}
		}

		if leader == "" {
			err = DeleteGame(gid)
			if err != nil {
				log.Printf("[Error] Could not delete game after leader left: %v", err)
			}
			return nil
		} else {
			game.Leader = leader
		}
	}

	err = db.UpdateGame(game)
	if err != nil {
		return fmt.Errorf("could not update game: %v", err)
	}

	return nil
}
