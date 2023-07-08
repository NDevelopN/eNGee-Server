package game

import u "Engee-Server/utils"

func CreateGame(g u.Game) (string, error) {

	return g.GID, nil
}

func GetGames() ([]u.Game, error) {
	return []u.Game{}, nil
}

func GetGame(gid string) (u.Game, error) {
	return u.Game{}, nil
}

func GetGamePlayers(gid string) ([]u.User, error) {
	return []u.User{}, nil
}

func UpdateGame(g u.Game) error {
	return nil
}

func ChangePlayerCount(g u.Game, d int) error {
	return nil
}

func DeleteGame(gid string) error {
	return nil
}
