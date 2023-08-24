package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"Engee-Server/utils"
)

var DB *sql.DB

/**
* This function initializes the database connection
* It checks for existing tables and currently removes them
* It then creates the games and players tables
 */
func InitDB() {
	var err error
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
		log.Fatalf("[Error] Failed to drop players table: %v", err)
	}
	_, err = DB.Query("DROP TABLE IF EXISTS gametypes;")
	if err != nil {
		log.Fatalf("[Error] Failed to drop gametypes table: %v", err)
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
		"uid 		varchar(80), " +
		"gid 		varchar(80), " +
		"name 		varchar(80), " +
		"status 	varchar(80) " +
		");")
	if err != nil {
		log.Fatalf("[Error] Failed to create players table: %v", err)
	}

	//TODO include handler
	_, err = DB.Query("CREATE TABLE gametypes ( " +
		"type 		varchar(80) " +
		");")
	if err != nil {
		log.Fatalf("[Error] Failed to create gametypes table: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("[Error] Failed to ping database: %v", err)
	}
}

/**
* This function returns all the games in the games table
 */
func GetAllGames() ([]utils.Game, error) {
	rows, err := DB.Query("SELECT * FROM games")
	if err != nil {
		return nil, fmt.Errorf("db failed: SELECT FROM games: %v", err)
	}
	defer rows.Close()

	gms := make([]utils.Game, 0)
	for rows.Next() {
		gm := new(utils.Game)
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
			return nil, fmt.Errorf("db failed: Scanning row: %v", err)
		}
		gms = append(gms, *gm)
	}

	err = rows.Err()

	if err == sql.ErrNoRows {
		return gms, nil
	} else if err != nil {
		return nil, fmt.Errorf("db failed: row error: %v", err)
	}

	return gms, nil
}

/**
* This function returns a single game, identified by the unique gid
 */
func GetGame(gid string) (utils.Game, error) {
	row := DB.QueryRow("SELECT * FROM games WHERE gid = $1", gid)

	gm := new(utils.Game)
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
		return *gm, err
	} else if err != nil {
		return *gm, fmt.Errorf("db failed: %v", err)
	}

	return *gm, nil
}

/**
 * This function adds a single new row to the games table
 */
func CreateGame(gm utils.Game) error {
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
		return fmt.Errorf("db failed: INSERT INTO games: %v", err)
	}

	return nil
}

/**
* This function updates a row in the games table, identified by the given game's gid
* All fields are updated in this function
 */
func UpdateGame(gm utils.Game) error {
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
	result, err := DB.Exec(
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
		return fmt.Errorf("db failed: UPDATE games: %v", err)
	}

	ra, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("db failed: failed to get rows affected: %v", err)
	}

	if ra <= 0 {
		return fmt.Errorf("db failed: no rows affected")
	}

	return nil
}

/**
* This function gets all players with the gid matching the given gid
* These are the players that are in the game
 */
func GetGamePlayers(gid string) ([]utils.User, error) {
	rows, err := DB.Query("SELECT * FROM players WHERE gid = $1", gid)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("db failed: SELECT FROM players: %v", err)
	}
	defer rows.Close()

	plrs := make([]utils.User, 0)
	for rows.Next() {
		plr := new(utils.User)
		err := rows.Scan(&plr.UID, &plr.GID, &plr.Name, &plr.Status)
		if err != nil {
			return nil, fmt.Errorf("db failed: scanning row: %v", err)
		}
		plrs = append(plrs, *plr)
	}
	if err == sql.ErrNoRows {
		return plrs, nil
	} else if err != nil {
		return nil, fmt.Errorf("db failed: row error")
	}

	return plrs, nil
}

/**
* This function deletes a row on the games table, identified by the given gid
 */
func RemoveGame(gid string) error {
	result, err := DB.Exec("DELETE FROM games WHERE gid = $1", gid)

	if err != nil {
		return fmt.Errorf("db failed: DELETE FROM games: %v", err)
	}

	ra, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("db failed: row error: %v", err)
	}

	if ra <= 0 {
		return fmt.Errorf("db failed: no rows affected")
	}

	return nil
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
* This function returns the player with the given uid
 */
func GetUser(uid string) (utils.User, error) {

	row := DB.QueryRow("SELECT * FROM players WHERE uid = $1", uid)

	plr := new(utils.User)
	err := row.Scan(&plr.UID, &plr.GID, &plr.Name, &plr.Status)
	if err == sql.ErrNoRows {
		return *plr, err
	} else if err != nil {
		return *plr, fmt.Errorf("db failed: %v", err)
	}

	return *plr, nil
}

/**
 * This function adds a single new row to the players table
 */
func CreateUser(user utils.User) error {
	_, err := DB.Exec("INSERT INTO players VALUES($1, $2, $3, $4)",
		user.UID,
		user.GID,
		user.Name,
		user.Status,
	)
	if err != nil {
		return fmt.Errorf("db failed: INSERT INTO players: %v", err)
	}

	return nil
}

/**
 * This function updates a row in the players table, identified by the given player's uid
 * All fields are updated in this function
 */
func UpdateUser(user utils.User) error {
	result, err := DB.Exec(
		"UPDATE players "+
			"SET "+
			"gid = $2, "+
			"name = $3, "+
			"status = $4 "+
			"WHERE uid = $1",
		user.UID, user.GID, user.Name, user.Status,
	)

	if err != nil {
		return fmt.Errorf("db failed: UPDATE players: %v", err)
	}

	ra, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("db failed: row error: %v", err)
	}

	if ra <= 0 {
		return fmt.Errorf("db failed: no rows affected")
	}

	return nil
}

/**
 * This function updates the status field of all players with gids that match the given gid
 */
func UpdateGameUserStatus(gid string, status string) error {

	_, err := DB.Exec(
		"UPDATE players "+
			"SET "+
			"status = $2 "+
			"WHERE gid = $1",
		gid, status)

	return fmt.Errorf("db failed: UPDATE players: %v", err)
}

/**
 * This function deletes a row on the palers table, identified by the given uid
 */
func RemoveUser(uid string) error {
	result, err := DB.Exec("DELETE FROM players WHERE uid = $1", uid)

	if err != nil {
		return fmt.Errorf("db failed: DELETE FROM players: %v", err)
	}

	ra, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("db failed: row error: %v", err)
	}

	if ra <= 0 {
		return fmt.Errorf("db failed: no rows affected")
	}

	return nil
}

func CreateGameTypes(gTypes map[string]utils.GHandler) error {
	//TODO include the handler
	for t := range gTypes {
		_, err := DB.Exec(`INSERT INTO gametypes (type) VALUES ($1)`, t)
		if err != nil {
			return fmt.Errorf("db failed: INSERT INTO types: %v", err)
		}
	}
	return nil
}

func GetGameTypes() ([]string, error) {
	rows, err := DB.Query("SELECT * FROM gametypes")
	if err != nil {
		return nil, fmt.Errorf("db failed: SELECT FROM gametypes: %v", err)
	}
	defer rows.Close()

	//TODO include handler
	typeList := []string{}
	for rows.Next() {
		var t string

		err = rows.Scan(&t)
		if err != nil {
			return nil, fmt.Errorf("db failed: Scanning row: %v", err)
		}

		typeList = append(typeList, t)
	}

	err = rows.Err()

	if err == sql.ErrNoRows {
		return typeList, nil
	} else if err != nil {
		return nil, fmt.Errorf("db failed: row error: %v", err)
	}

	return typeList, nil
}
