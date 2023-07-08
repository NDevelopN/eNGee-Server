package game

import (
	"testing"

	db "Engee-Server/database"
	u "Engee-Server/utils"

	"github.com/google/uuid"
)

// Test Game Creation
func TestCreateGameValid(t *testing.T) {
	db.InitDB()
	msg, err := CreateGame(u.DefGame)
	_, pe := uuid.Parse(msg)
	if pe != nil || err != nil {
		t.Fatalf(`CreateGame(valid) = %q, %q, want "uuid", "nil"`, msg, err)
	}
}

func TestCreateGameMulti(t *testing.T) {
	db.InitDB()
	_, _ = CreateGame(u.DefGame)
	msg, err := CreateGame(u.DefGame)
	_, pe := uuid.Parse(msg)
	if pe != nil || err != nil {
		t.Fatalf(`CreateGame(multi-valid) = %q, %q, want "uuid", "nil"`, msg, err)
	}
}
func TestCreateGameEmptyValues(t *testing.T) {
	db.InitDB()
	game := u.DefGame
	game.Name = ""
	game.Type = ""

	msg, err := CreateGame(game)
	if msg != "" || err == nil {
		t.Fatalf(`CreateGame(EmptyVals) = %q, %q, want "", error`, msg, err)
	}
}

func TestCreateGameInvalidPlrNums(t *testing.T) {
	db.InitDB()

	game := u.DefGame
	game.MinPlrs = 10

	msg, err := CreateGame(game)
	if msg != "" || err == nil {
		t.Fatalf(`CreateGame(InvalidPlrNums) = %q, %q, want "", error`, msg, err)
	}
}

func TestCreateGameInjection(t *testing.T) {
	db.InitDB()
	//TODO
}

// Test Game Retrieval
func TestGetGamesValidSingle(t *testing.T) {
	db.InitDB()
	gid, _ := CreateGame(u.DefGame)

	want := u.DefGame
	want.GID = gid

	games, err := GetGames()
	if want != games[0] || err != nil {
		t.Fatalf(`GetGames(Valid) = %q, %q, want %q, "nil"`, games, err, want)
	}
}

func TestGetGamesValidMulti(t *testing.T) {
	db.InitDB()
	gid1, _ := CreateGame(u.DefGame)
	gid2, _ := CreateGame(u.DefGame)

	var want [2]u.Game
	want[0] = u.DefGame
	want[0].GID = gid1
	want[1] = u.DefGame
	want[1].GID = gid2

	games, err := GetGames()

	if want[0] != games[0] || want[1] != games[1] || err != nil {
		t.Fatalf(`GetGames(ValidMulti) = %q, %q, want %q, "nil"`, games, err, want)
	}
}

func TestGetGamesEmpty(t *testing.T) {
	db.InitDB()
	games, err := GetGames()

	if len(games) > 0 || err != nil {
		t.Fatalf(`GetGames(Empty) = %q, %q, want "[]", "nil"`, games, err)
	}
}

// Test Get Game
func TestGetGameValid(t *testing.T) {
	db.InitDB()

	gid, _ := CreateGame(u.DefGame)

	game, err := GetGame(gid)
	if err != nil {
		t.Fatalf(`GetGame(valid) = Failed to get game: %q`, err)
	}

	game.GID = ""

	if game != u.DefGame {
		t.Fatalf(`GetGame(valid) = %q, want %q`, game, u.DefGame)
	}
}

func TestGetGameMulti(t *testing.T) {
	db.InitDB()

	_, _ = CreateGame(u.DefGame)
	gid, _ := CreateGame(u.DefGame)

	game, err := GetGame(gid)
	if err != nil {
		t.Fatalf(`GetGame(valid) = Failed to get game: %q`, err)
	}

	game.GID = ""

	if game != u.DefGame {
		t.Fatalf(`GetGame(valid) = %q, want %q`, game, u.DefGame)
	}
}

func TestGetGameInvalidGID(t *testing.T) {
	db.InitDB()

	_, _ = CreateGame(u.DefGame)

	_, err := GetGame(uuid.NewString())
	if err == nil {
		t.Fatalf(`GetGame(InvalidGID) = %q, want ERROR`, err)
	}
}

func TestGetGameEmptyGID(t *testing.T) {
	db.InitDB()

	_, _ = CreateGame(u.DefGame)

	_, err := GetGame("")
	if err == nil {
		t.Fatalf(`GetGame(EmptyGID) = %q, want ERROR`, err)
	}
}

func TestGetGameEmptyDB(t *testing.T) {
	db.InitDB()

	_, err := GetGame(uuid.NewString())
	if err == nil {
		t.Fatalf(`GetGame(EmptyDB) = %q, want ERROR`, err)
	}
}

func TestGetGameInjection(t *testing.T) {
	db.InitDB()
	//TODO
}

// Test Game Update
func TestUpdateGameChangeName(t *testing.T) {
	db.InitDB()

	gid, _ := CreateGame(u.DefGame)

	game := u.DefGame
	game.GID = gid
	game.Name = "Game Test"

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(Name) = %q, want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(Name) = %q, %q, want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeType(t *testing.T) {
	db.InitDB()

	gid, _ := CreateGame(u.DefGame)

	game := u.DefGame
	game.GID = gid
	game.Type = "TypeTest"

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(Type) = %q, want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(Type) = %q, %q, want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeStatus(t *testing.T) {
	db.InitDB()

	gid, _ := CreateGame(u.DefGame)

	game := u.DefGame
	game.GID = gid
	game.Status = "Test Status"

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(Status) = %q, want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(Status) = %q, %q, want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeOldStatus(t *testing.T) {
	db.InitDB()

	gid, _ := CreateGame(u.DefGame)

	game := u.DefGame
	game.GID = gid
	game.OldStatus = "Old Test Status"

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(OldStatus) = %q, want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(OldStatus) = %q, %q, want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeLeader(t *testing.T) {
	db.InitDB()

	gid, _ := CreateGame(u.DefGame)

	game := u.DefGame
	game.GID = gid
	game.Leader = uuid.NewString()

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(Leader) = %q, want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(Leader) = %q, %q, want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeMinPlrs(t *testing.T) {
	db.InitDB()

	gid, _ := CreateGame(u.DefGame)

	game := u.DefGame
	game.GID = gid
	game.MinPlrs = 2

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(MinP) = %q, want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(MinP) = %q, %q, want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeMaxPlrs(t *testing.T) {
	db.InitDB()

	gid, _ := CreateGame(u.DefGame)

	game := u.DefGame
	game.GID = gid
	game.MaxPlrs = 10

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(MaxP) = %q, want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(MaxP) = %q, %q, want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeMinPlrsHigh(t *testing.T) {
	db.InitDB()

	gid, _ := CreateGame(u.DefGame)

	game := u.DefGame
	game.GID = gid
	game.MinPlrs = 10

	err := UpdateGame(game)
	if err == nil {
		t.Fatalf(`UpdateGame(MinHigh) = %q, want ERROR`, err)
	}

	game = u.DefGame
	game.GID = gid

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(MinHigh) = %q, %q, want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeMaxPlrsLow(t *testing.T) {
	db.InitDB()
	gid, _ := CreateGame(u.DefGame)

	game := u.DefGame
	game.GID = gid
	game.MaxPlrs = 0

	err := UpdateGame(game)
	if err == nil {
		t.Fatalf(`UpdateGame(MaxLow) = %q, want ERROR`, err)
	}

	game = u.DefGame
	game.GID = gid

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(MaxLow) = %q, %q, want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeCurPlrs(t *testing.T) {
	db.InitDB()

	gid, _ := CreateGame(u.DefGame)

	game := u.DefGame
	game.GID = gid
	game.CurPlrs = 1

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(CurP) = %q, want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(CurP) = %q, %q, want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeCurHigh(t *testing.T) {
	db.InitDB()

	gid, _ := CreateGame(u.DefGame)

	game := u.DefGame
	game.GID = gid
	game.CurPlrs = 99

	err := UpdateGame(game)
	if err == nil {
		t.Fatalf(`UpdateGame(CurPHigh) = %q, want ERROR`, err)
	}

	game = u.DefGame
	game.GID = gid

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(CurPHigh) = %q, %q, want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeAdditional(t *testing.T) {
	db.InitDB()

	gid, _ := CreateGame(u.DefGame)

	game := u.DefGame
	game.GID = gid
	game.AdditionalRules = `{"rule1": "default"}`

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(AddRules) = %q, want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(AddRules) = %q, %q, want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeAll(t *testing.T) {
	db.InitDB()

	gid, _ := CreateGame(u.DefGame)

	var game = u.Game{
		GID:             gid,
		Name:            "Game Test",
		Type:            "TypeTest",
		Status:          "Test Status",
		OldStatus:       "Old Test Status",
		Leader:          uuid.NewString(),
		MinPlrs:         2,
		MaxPlrs:         10,
		CurPlrs:         1,
		AdditionalRules: `{"rule1": "default"}`,
	}

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(All) = %q, want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(All) = %q, %q, want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameInvalidGID(t *testing.T) {
	db.InitDB()

	gid, _ := CreateGame(u.DefGame)

	game := u.DefGame
	game.GID = uuid.NewString()
	game.Name = "Game Test"

	err := UpdateGame(game)
	if err == nil {
		t.Fatalf(`UpdateGame(InvalidGID) = %q, want ERROR`, err)
	}

	game = u.DefGame
	game.GID = gid

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(InvalidGID) = %q, %q, want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameEmptyGID(t *testing.T) {
	db.InitDB()

	gid, _ := CreateGame(u.DefGame)

	game := u.DefGame
	game.GID = ""
	game.Name = "Game Test"

	err := UpdateGame(game)
	if err == nil {
		t.Fatalf(`UpdateGame(EmptyGID) = %q, want ERROR`, err)
	}

	game = u.DefGame
	game.GID = gid

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(EmptyGID) = %q, %q, want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameEmptyDB(t *testing.T) {
	db.InitDB()

	game := u.DefGame
	game.GID = uuid.NewString()
	game.Name = "Game Test"

	err := UpdateGame(game)
	if err == nil {
		t.Fatalf(`UpdateGame(EmptyDB) = %q, want ERROR`, err)
	}

	games, err := GetGames()
	if len(games) > 0 || err != nil {
		t.Fatalf(`UpdateGame(EmptyDB) - GET = %q, %q, want "[]", "nil"`, games, err)
	}
}

func TestUpdateGameNoChange(t *testing.T) {
	db.InitDB()

	gid, _ := CreateGame(u.DefGame)

	game := u.DefGame
	game.GID = gid

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(NoChange) = %q, want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(NoChange) = %q, %q, want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameInjection(t *testing.T) {
	db.InitDB()
	// TODO
}

func TestChangePlayerCountValidIncrease(t *testing.T) {
	db.InitDB()

	game := u.DefGame
	game.CurPlrs = 2

	gid, _ := CreateGame(game)
	game.GID = gid

	err := ChangePlayerCount(game, 1)
	if err != nil {
		t.Fatalf(`ChangePlayerCount(ValidIncrease) = %q, want "nil"`, err)
	}

	want := game.CurPlrs + 1

	game, err = GetGame(gid)
	if game.CurPlrs != want || err != nil {
		t.Fatalf(`ChangePlayerCount(ValidIncrease) = %q, %q, want %q, "nil"`, game.CurPlrs, err, want)
	}
}

func TestChangePlayerCountValidDecrease(t *testing.T) {
	db.InitDB()

	game := u.DefGame
	game.CurPlrs = 2

	gid, _ := CreateGame(u.DefGame)

	game.GID = gid

	err := ChangePlayerCount(game, -1)
	if err != nil {
		t.Fatalf(`ChangePlayerCount(ValidDecrease) = %q, want "nil"`, err)
	}

	want := game.CurPlrs - 1

	game, err = GetGame(gid)
	if game.CurPlrs != want || err != nil {
		t.Fatalf(`ChangePlayerCount(ValidDecrease) = %q, %q, want %q, "nil"`, game.CurPlrs, err, want)
	}
}

func TestChangePlayerCountDoubleIncrease(t *testing.T) {
	db.InitDB()

	game := u.DefGame
	game.CurPlrs = 1

	gid, _ := CreateGame(u.DefGame)
	game.GID = gid

	_ = ChangePlayerCount(game, 1)
	err := ChangePlayerCount(game, 1)
	if err != nil {
		t.Fatalf(`ChangePlayerCount(DoubleIncrease) = %q, want "nil"`, err)
	}

	want := game.CurPlrs + 2

	game, err = GetGame(gid)
	if game.CurPlrs != want || err != nil {
		t.Fatalf(`ChangePlayerCount(DoubleIncrease) = %q, %q, want %q, "nil"`, game.CurPlrs, err, want)
	}
}

func TestChangePlayerCountDoubleDecrease(t *testing.T) {
	db.InitDB()

	game := u.DefGame
	game.CurPlrs = 3

	gid, _ := CreateGame(u.DefGame)
	game.GID = gid

	_ = ChangePlayerCount(game, -1)
	err := ChangePlayerCount(game, -1)
	if err != nil {
		t.Fatalf(`ChangePlayerCount(DoubleDecrease) = %q, want "nil"`, err)
	}

	want := game.CurPlrs - 2

	game, err = GetGame(gid)
	if game.CurPlrs != want || err != nil {
		t.Fatalf(`ChangePlayerCount(DoubleDecrease) = %q, %q, want %q, "nil"`, game.CurPlrs, err, want)
	}
}
func TestChangePlayerCountIncreaseFullGame(t *testing.T) {
	db.InitDB()

	game := u.DefGame
	game.CurPlrs = game.MaxPlrs

	gid, _ := CreateGame(game)

	game.GID = gid

	err := ChangePlayerCount(game, 1)
	if err == nil {
		t.Fatalf(`ChangePlayerCount(FullGameIncrease) = %q, want ERROR`, err)
	}

	want := game.CurPlrs

	game, err = GetGame(gid)
	if game.CurPlrs != want || err != nil {
		t.Fatalf(`ChangePlayerCount(FullGameIncrease) = %q, %q, want %q, "nil"`, game.CurPlrs, err, want)
	}
}

func TestChangePlayerCountDecreaseToZero(t *testing.T) {
	db.InitDB()

	game := u.DefGame
	game.CurPlrs = 1

	gid, _ := CreateGame(game)

	game.GID = gid

	err := ChangePlayerCount(game, -1)
	if err != nil {
		t.Fatalf(`ChangePlayerCount(DecreaseToZero) = %q, want "nil"`, err)
	}

	game, err = GetGame(gid)
	if err == nil {
		t.Fatalf(`ChangePlayerCount(DecreaseToZero) = %q, %q, want "nil", ERROR`, game, err)
	}
}

func TestChangePlayerCountNoChange(t *testing.T) {
	db.InitDB()

	game := u.DefGame
	game.CurPlrs = 1

	gid, _ := CreateGame(game)

	game.GID = gid

	err := ChangePlayerCount(game, 0)
	if err != nil {
		t.Fatalf(`ChangePlayerCount(NoChange) = %q, want "nil"`, err)
	}

	want := game.CurPlrs

	game, err = GetGame(gid)
	if game.CurPlrs != want || err != nil {
		t.Fatalf(`ChangePlayerCount(NoChange) %q, %q, want %q, "nil"`, game.CurPlrs, err, want)
	}
}

// Test Game Deletion
func TestDeleteGameValid(t *testing.T) {
	db.InitDB()
	gid, _ := CreateGame(u.DefGame)
	err := DeleteGame(gid)
	if err != nil {
		t.Fatalf(`DeleteGame(Valid) = %q, want "nil"`, err)
	}

	games, err := GetGames()
	if len(games) > 0 || err != nil {
		t.Fatalf(`DeleteGame(Valid) = %q, %q, want 0, "nil"`, len(games), err)
	}
}

func TestDeleteGameMulti(t *testing.T) {
	db.InitDB()
	_, _ = CreateGame(u.DefGame)
	gid, _ := CreateGame(u.DefGame)
	err := DeleteGame(gid)
	if err != nil {
		t.Fatalf(`DeleteGame(Multi) = %q, want  "nil"`, err)
	}

	games, err := GetGames()
	if len(games) > 1 || err != nil {
		t.Fatalf(`DeletGame(Multi) = %q, %q, want 1, "nil"`, len(games), err)
	}
}

func TestDeleteGameInvalidGID(t *testing.T) {
	db.InitDB()
	_, _ = CreateGame(u.DefGame)

	err := DeleteGame(uuid.NewString())
	if err == nil {
		t.Fatalf(`DeleteGame(InvalidGID) = %q, want ERROR`, err)
	}

	games, err := GetGames()
	if len(games) < 1 || err != nil {
		t.Fatalf(`DeletGame(InvalidGID) = %q, %q, want 0, "nil"`, len(games), err)
	}
}

func TestDeletGameEmptyGID(t *testing.T) {
	db.InitDB()
	_, _ = CreateGame(u.DefGame)

	err := DeleteGame("")
	if err == nil {
		t.Fatalf(`DeleteGame(EmptyGID) = %q, want ERROR`, err)
	}

	games, err := GetGames()
	if len(games) < 1 || err != nil {
		t.Fatalf(`DeletGame(EmptyGID) = %q, %q, want 0, "nil"`, len(games), err)
	}
}

func TestDeleteGameEmptyDB(t *testing.T) {
	db.InitDB()

	err := DeleteGame(uuid.NewString())
	if err == nil {
		t.Fatalf(`DeleteGame(EmptyDB) = %q, want ERROR`, err)
	}
}

func TestDeleteGameRepeat(t *testing.T) {
	db.InitDB()
	gid, _ := CreateGame(u.DefGame)
	_ = DeleteGame(gid)

	err := DeleteGame(gid)
	if err == nil {
		t.Fatalf(`DeleteGame(Repeat) = %q, want ERROR`, err)
	}
}

func TestDeleteGameInjection(t *testing.T) {
	db.InitDB()
	//TODO

}
