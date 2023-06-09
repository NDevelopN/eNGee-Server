package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	u "Engee-Server/utils"
)

var DB *sql.DB

/**
* This function initializes the database connection
* It checks for existing tables and currently removes them
* It then creates the games and players tables
* TODO: An external file with more permanent storage of tables
 */
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

/**
* This function returns all the games in the games table
 */
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

/**
* This function returns a single game, identified by the unique gid
 */
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

/**
 * This function adds a single new row to the games table
 */
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

/**
* This function updates a row in the games table, identified by the given game's gid
* All fields are updated in this function
 */
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

/**
* This function gets all players with the gid matching the given gid
* These are the players that are in the game
 */
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

/**
* This function deletes a row on the games table, identified by the given gid
 */
func RemoveGame(gid string) error {
	_, err := DB.Query("DELETE FROM games WHERE gid = $1", gid)
	return err
}

/**
* This function returns a count of players who have matching gids with the given gid
 */
func GetGamePCount(gid string) int {
	count := 0
	DB.QueryRow("SELECT count(*) FROM players WHERE gid = $1", gid).Scan(&count)

	return count
}

/**
* This function returns a count of players who have the correct gid and a status of Ready
 */
func GetGamePReady(gid string) int {
	count := 0
	DB.QueryRow("SELECT count(*) FROM players WHERE gid = $1 AND status = $2", gid, "Ready").Scan(&count)

	return count
}

/**
* This function returns the player with the given pid
 */
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

/**
 * This function adds a single new row to the players table
 */
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

/**
 * This function updates a row in the players table, identified by the given player's pid
 * All fields are updated in this function
 */
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

/**
 * This function updates the status field of all players with gids that match the given gid
 */
func UpdateGamePlayerStatus(gid string, status string) error {

	_, err := DB.Exec(
		"UPDATE players "+
			"SET "+
			"status = $2 "+
			"WHERE gid = $1",
		gid, status)

	return err

}

/**
 * This function deletes a row on the palers table, identified by the given pid
 */
func RemovePlayer(pid string) error {
	_, err := DB.Query("DELETE FROM players WHERE pid = $1", pid)
	return err
}
