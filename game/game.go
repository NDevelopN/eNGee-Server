package game

import (
	db "Engee-Server/database"
	"Engee-Server/utils"
	"fmt"

	"github.com/google/uuid"
)

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

	utils.AddConnectionPool(g.GID)

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
	if err != nil {
		return nil, fmt.Errorf("failed to find game in database: %v", err)
	}

	plrs, err := db.GetGamePlayers(gid)
	if err != nil {
		return nil, fmt.Errorf("failed to get players from database: %v", err)
	}

	return plrs, nil
}

func UpdateGame(g utils.Game) error {
	if g.Name == "" {
		return fmt.Errorf("cannot set game name to empty string")
	}

	if g.Type == "" {
		return fmt.Errorf("cannot set game type to empty string")
	}

	if g.Status == "" {
		return fmt.Errorf("cannot set game status to empty string")
	}

	if g.MinPlrs > g.MaxPlrs {
		return fmt.Errorf("provided minPlrs is greater than provided maxPlrs")
	}

	if g.CurPlrs > g.MaxPlrs {
		return fmt.Errorf("provided curPlrs is greater than provided maxPlrs")
	}

	err := db.UpdateGame(g)
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
		return fmt.Errorf("could not get players from database: %v", err)
	}

	for _, p := range plrs {
		p.GID = ""
		err = db.UpdateUser(p)
		if err != nil {
			return fmt.Errorf("could not update player (clearing GID) in database: %v", err)
		}
	}

	utils.RemoveConnectionPool(gid)

	err = db.RemoveGame(gid)
	if err != nil {
		return fmt.Errorf("could not delete the game from database: %v", err)
	}

	return nil
}
