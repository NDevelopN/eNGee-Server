package game

import (
	db "Engee-Server/database"
	u "Engee-Server/utils"
	"fmt"

	"github.com/google/uuid"
)

func CreateGame(g u.Game) (string, error) {
	if g.Name == "" {
		return "", fmt.Errorf("provided game name is empty")
	}

	if g.Type == "" {
		return "", fmt.Errorf("provided game type is empty")
	}

	if g.MinPlrs > g.MaxPlrs {
		return "", fmt.Errorf("provided minPlrs is greater than provided maxPlrs")
	}

	g.GID = uuid.NewString()
	g.Status = "Lobby"

	err := db.CreateGame(g)
	if err != nil {
		return "", fmt.Errorf("failed to create game in database: %v", err)
	}

	return g.GID, nil
}

func GetGames() ([]u.Game, error) {
	//TODO any checks needed here?
	return db.GetAllGames()
}

func GetGame(gid string) (u.Game, error) {
	//TODO any checks needed here?
	return db.GetGame(gid)
}

func GetGamePlayers(gid string) ([]u.User, error) {
	return []u.User{}, nil
}

func UpdateGame(g u.Game) error {
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

func ChangePlayerCount(g u.Game, d int) error {
	return nil
}

func DeleteGame(gid string) error {
	return nil
}
