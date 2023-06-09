package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	u "Engee-Server/utils"
)

var DB *sql.DB

func InitDB() {
	var err error
	//TODO config file
	DB, err = sql.Open("postgres", "postgres://ngdbu:ngp@localhost/ngdb")
	if err != nil {
		log.Fatalf("[Error] Failed to open connection to sql server: %v", err)
	}

	_, err = DB.Query("DROP TABLE IF EXISTS games;")
	if err != nil {
		log.Fatalf("[Error] Failed to drop games table: %v", err)
	}
	_, err = DB.Query("DROP TABLE IF EXISTS players;")
	if err != nil {
		log.Fatalf("[Error] Failed to drop games table: %v", err)
	}

	_, err = DB.Query("CREATE TABLE games (" +
		"gid 		varchar(80), " +
		"name 		varchar(80), " +
		"type 		varchar(80), " +
		"status 	varchar(80), " +
		"ostatus 	varchar(80), " +
		"leader 	varchar(80), " +
		"min_plrs 	smallint, " +
		"max_plrs 	smallint, " +
		"cur_plrs 	smallint, " +
		"add_rules 	varchar(255) " +
		");")
	if err != nil {
		log.Fatalf("[Error] Failed to create games table: %v", err)
	}

	_, err = DB.Query("CREATE TABLE players ( " +
		"pid 		varchar(80), " +
		"gid 		varchar(80), " +
		"name 		varchar(80), " +
		"status 	varchar(80) " +
		");")
	if err != nil {
		log.Fatalf("[Error] Failed to create players table: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("[Error] Failed to ping database: %v", err)
	}
}

func GetAllGames() ([]u.Game, error) {
	rows, err := DB.Query("SELECT * FROM games")
	if err != nil {
		log.Printf("[Error] Failed to get games from db: %v", err)
		return nil, err
	}
	defer rows.Close()

	gms := make([]u.Game, 0)
	for rows.Next() {
		gm := new(u.Game)
		err := rows.Scan(
			&gm.GID,
			&gm.Name,
			&gm.Type,
			&gm.Status,
			&gm.OldStatus,
			&gm.Leader,
			&gm.MinPlrs,
			&gm.MaxPlrs,
			&gm.CurPlrs,
			&gm.AdditionalRules,
		)
		if err != nil {
			log.Printf("[Error] Failed to Scan row in games: %v", err)
			return nil, err
		}
		gms = append(gms, *gm)
	}

	if err = rows.Err(); err != nil {
		log.Printf("[Error] Error while scanning row in games: %v", err)
		return nil, err
	}

	return gms, nil
}

func GetGame(gid string) (u.Game, error) {
	row := DB.QueryRow("SELECT * FROM games WHERE gid = $1", gid)

	gm := new(u.Game)
	err := row.Scan(
		&gm.GID,
		&gm.Name,
		&gm.Type,
		&gm.Status,
		&gm.OldStatus,
		&gm.Leader,
		&gm.MinPlrs,
		&gm.MaxPlrs,
		&gm.CurPlrs,
		&gm.AdditionalRules,
	)
	if err == sql.ErrNoRows {
		log.Printf("[Warn] Did not find Game %v", err)
		return *gm, err
	} else if err != nil {
		log.Printf("[Error] while searching for Game %v", err)
		return *gm, err
	}

	return *gm, nil
}

func CreateGame(gm u.Game) error {
	createStatement := `
		INSERT INTO games (
			gid,
			name,
			type, 
			status,
			ostatus,
			leader,
			min_plrs,
			max_plrs,
			cur_plrs,
			add_rules
		) Values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := DB.Exec(
		createStatement,
		gm.GID,
		gm.Name,
		gm.Type,
		gm.Status,
		gm.OldStatus,
		gm.Leader,
		int(gm.MinPlrs),
		int(gm.MaxPlrs),
		int(gm.CurPlrs),
		gm.AdditionalRules,
	)
	if err != nil {
		log.Printf("[Error] Could not create game: %v", err)
		return err
	}

	return nil

}

func UpdateGame(gm u.Game) error {
	updateStatement := `
		UPDATE games
		SET name = $2, 
			type = $3, 
			status = $4, 
			ostatus = $5, 
			leader = $6, 
			min_plrs = $7, 
			max_plrs = $8, 
			cur_plrs = $9, 
			add_rules = $10
		WHERE gid = $1;`
	_, err := DB.Exec(
		updateStatement,
		gm.GID,
		gm.Name,
		gm.Type,
		gm.Status,
		gm.OldStatus,
		gm.Leader,
		int(gm.MinPlrs),
		int(gm.MaxPlrs),
		int(gm.CurPlrs),
		gm.AdditionalRules,
	)
	if err != nil {
		log.Printf("[Error] Failed to update game: %v", err)
	}

	return err
}

func GetGamePlayers(gid string) ([]u.Player, error) {
	rows, err := DB.Query("SELECT * FROM players WHERE gid = $1", gid)
	if err != nil {
		log.Printf("[Error] Failed to get game players from database: %v", err)
		return nil, err
	}
	defer rows.Close()

	plrs := make([]u.Player, 0)
	for rows.Next() {
		plr := new(u.Player)
		err := rows.Scan(&plr.PID, &plr.GID, &plr.Name, &plr.Status)
		if err != nil {
			log.Printf("[Error] Failed to Scan row in game players: %v", err)
			return nil, err
		}
		plrs = append(plrs, *plr)
	}
	if err = rows.Err(); err != nil {
		log.Printf("[Error] while scanning rows in game players: %v", err)
		return nil, err
	}

	return plrs, nil
}

func RemoveGame(gid string) error {
	_, err := DB.Query("DELETE FROM games WHERE gid = $1", gid)
	return err
}

func GetGamePCount(gid string) int {
	count := 0
	DB.QueryRow("SELECT count(*) FROM players WHERE gid = $1", gid).Scan(&count)

	return count
}

func GetGamePReady(gid string) int {
	count := 0
	DB.QueryRow("SELECT count(*) FROM players WHERE gid = $1 AND status = $2", gid, "Ready").Scan(&count)

	return count
}

func GetPlayer(pid string) (u.Player, error) {

	row := DB.QueryRow("SELECT * FROM players WHERE pid = $1", pid)

	plr := new(u.Player)
	err := row.Scan(&plr.PID, &plr.GID, &plr.Name, &plr.Status)
	if err == sql.ErrNoRows {
		log.Printf("[Error] Did not find Player %v", err)
		return *plr, err
	} else if err != nil {
		log.Printf("[Error] Error searching for player %v", err)
		return *plr, err
	}

	return *plr, nil
}

func CreatePlayer(plr u.Player) error {
	_, err := DB.Exec("INSERT INTO players VALUES($1, $2, $3, $4)",
		plr.PID,
		plr.GID,
		plr.Name,
		plr.Status,
	)
	if err != nil {
		log.Printf("[Error] Could not create player: %v", err)
		return err
	}

	return nil
}

func UpdatePlayer(plr u.Player) error {
	_, err := DB.Exec(
		"UPDATE players "+
			"SET "+
			"gid = $2, "+
			"name = $3, "+
			"status = $4 "+
			"WHERE pid = $1",
		plr.PID, plr.GID, plr.Name, plr.Status,
	)

	return err
}

func UpdateGamePlayerStatus(gm u.Game, status string) error {

	_, err := DB.Exec(
		"UPDATE players "+
			"SET "+
			"status = $2 "+
			"WHERE gid = $1",
		gm.GID, status)

	return err

}

func RemovePlayer(pid string) error {
	_, err := DB.Query("DELETE FROM players WHERE pid = $1", pid)
	return err
}
