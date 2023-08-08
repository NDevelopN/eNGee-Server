package game

import (
	"testing"

	db "Engee-Server/database"
	"Engee-Server/utils"

	"github.com/google/uuid"
)

func prepTest() {
	db.InitDB()
	utils.NO_HANDLER = true
}

// Test Game Creation
func TestCreateGameValid(t *testing.T) {
	prepTest()
	msg, err := CreateGame(utils.DefGame)
	_, pe := uuid.Parse(msg)
	if pe != nil || err != nil {
		t.Fatalf(`CreateGame(valid) = %q, "%v", want "uuid", "nil"`, msg, err)
	}
}

func TestCreateGameMulti(t *testing.T) {
	prepTest()
	_, _ = CreateGame(utils.DefGame)
	msg, err := CreateGame(utils.DefGame)
	_, pe := uuid.Parse(msg)
	if pe != nil || err != nil {
		t.Fatalf(`CreateGame(multi-valid) = %q, "%v", want "uuid", "nil"`, msg, err)
	}
}
func TestCreateGameEmptyValues(t *testing.T) {
	prepTest()
	game := utils.DefGame
	game.Name = ""
	game.Type = ""

	msg, err := CreateGame(game)
	if msg != "" || err == nil {
		t.Fatalf(`CreateGame(EmptyVals) = %q, "%v", want "", error`, msg, err)
	}
}

func TestCreateGameInvalidPlrNums(t *testing.T) {
	prepTest()

	game := utils.DefGame
	game.MinPlrs = 10

	msg, err := CreateGame(game)
	if msg != "" || err == nil {
		t.Fatalf(`CreateGame(InvalidPlrNums) = %q, "%v", want "", error`, msg, err)
	}
}

func TestCreateGameInjection(t *testing.T) {
	prepTest()
	//TODO
}

// Test Game Retrieval
func TestGetGamesValidSingle(t *testing.T) {
	prepTest()
	gid, _ := CreateGame(utils.DefGame)

	want := utils.DefGame
	want.GID = gid

	games, err := GetGames()
	if want != games[0] || err != nil {
		t.Fatalf(`GetGames(Valid) = %q, "%v", want %q, "nil"`, games, err, want)
	}
}

func TestGetGamesValidMulti(t *testing.T) {
	prepTest()
	gid1, _ := CreateGame(utils.DefGame)
	gid2, _ := CreateGame(utils.DefGame)

	var want [2]utils.Game
	want[0] = utils.DefGame
	want[0].GID = gid1
	want[1] = utils.DefGame
	want[1].GID = gid2

	games, err := GetGames()

	if want[0] != games[0] || want[1] != games[1] || err != nil {
		t.Fatalf(`GetGames(ValidMulti) = %q, "%v", want %q, "nil"`, games, err, want)
	}
}

func TestGetGamesEmpty(t *testing.T) {
	prepTest()
	games, err := GetGames()

	if len(games) > 0 || err != nil {
		t.Fatalf(`GetGames(Empty) = %q, "%v", want "[]", "nil"`, games, err)
	}
}

// Test Get Game
func TestGetGameValid(t *testing.T) {
	prepTest()

	gid, _ := CreateGame(utils.DefGame)

	game, err := GetGame(gid)
	if err != nil {
		t.Fatalf(`GetGame(valid) = Failed to get game: "%v"`, err)
	}

	game.GID = ""

	if game != utils.DefGame {
		t.Fatalf(`GetGame(valid) = %q, want %q`, game, utils.DefGame)
	}
}

func TestGetGameMulti(t *testing.T) {
	prepTest()

	_, _ = CreateGame(utils.DefGame)
	gid, _ := CreateGame(utils.DefGame)

	game, err := GetGame(gid)
	if err != nil {
		t.Fatalf(`GetGame(valid) = Failed to get game: "%v"`, err)
	}

	game.GID = ""

	if game != utils.DefGame {
		t.Fatalf(`GetGame(valid) = %q, want %q`, game, utils.DefGame)
	}
}

func TestGetGameInvalidGID(t *testing.T) {
	prepTest()

	_, _ = CreateGame(utils.DefGame)

	_, err := GetGame(uuid.NewString())
	if err == nil {
		t.Fatalf(`GetGame(InvalidGID) = "%v", want ERROR`, err)
	}
}

func TestGetGameEmptyGID(t *testing.T) {
	prepTest()

	_, _ = CreateGame(utils.DefGame)

	_, err := GetGame("")
	if err == nil {
		t.Fatalf(`GetGame(EmptyGID) = "%v", want ERROR`, err)
	}
}

func TestGetGameEmptyDB(t *testing.T) {
	prepTest()

	_, err := GetGame(uuid.NewString())
	if err == nil {
		t.Fatalf(`GetGame(EmptyDB) = "%v", want ERROR`, err)
	}
}

func TestGetGameInjection(t *testing.T) {
	prepTest()
	//TODO
}

// Test Game Update
func TestUpdateGameChangeName(t *testing.T) {
	prepTest()

	gid, _ := CreateGame(utils.DefGame)

	game := utils.DefGame
	game.GID = gid
	game.Name = "Game Test"

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(Name) = "%v", want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(Name) = %q, "%v", want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeType(t *testing.T) {
	prepTest()

	gid, _ := CreateGame(utils.DefGame)

	game := utils.DefGame
	game.GID = gid
	game.Type = "TypeTest"

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(Type) = "%v", want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(Type) = %q, "%v", want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeStatus(t *testing.T) {
	prepTest()

	gid, _ := CreateGame(utils.DefGame)

	game := utils.DefGame
	game.GID = gid
	game.Status = "Test Status"

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(Status) = "%v", want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(Status) = %q, "%v", want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeOldStatus(t *testing.T) {
	prepTest()

	gid, _ := CreateGame(utils.DefGame)

	game := utils.DefGame
	game.GID = gid
	game.OldStatus = "Old Test Status"

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(OldStatus) = "%v", want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(OldStatus) = %q, "%v", want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeLeader(t *testing.T) {
	prepTest()

	gid, _ := CreateGame(utils.DefGame)

	game := utils.DefGame
	game.GID = gid
	game.Leader = uuid.NewString()

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(Leader) = "%v", want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(Leader) = %q, "%v", want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeMinPlrs(t *testing.T) {
	prepTest()

	gid, _ := CreateGame(utils.DefGame)

	game := utils.DefGame
	game.GID = gid
	game.MinPlrs = 2

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(MinP) = "%v", want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(MinP) = %q, "%v", want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeMaxPlrs(t *testing.T) {
	prepTest()

	gid, _ := CreateGame(utils.DefGame)

	game := utils.DefGame
	game.GID = gid
	game.MaxPlrs = 10

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(MaxP) = "%v", want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(MaxP) = %q, "%v", want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeMinPlrsHigh(t *testing.T) {
	prepTest()

	gid, _ := CreateGame(utils.DefGame)

	game := utils.DefGame
	game.GID = gid
	game.MinPlrs = 10

	err := UpdateGame(game)
	if err == nil {
		t.Fatalf(`UpdateGame(MinHigh) = "%v", want ERROR`, err)
	}

	game = utils.DefGame
	game.GID = gid

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(MinHigh) = %q, "%v", want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeMaxPlrsLow(t *testing.T) {
	prepTest()
	gid, _ := CreateGame(utils.DefGame)

	game := utils.DefGame
	game.GID = gid
	game.MaxPlrs = 0

	err := UpdateGame(game)
	if err == nil {
		t.Fatalf(`UpdateGame(MaxLow) = "%v", want ERROR`, err)
	}

	game = utils.DefGame
	game.GID = gid

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(MaxLow) = %q, "%v", want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeCurPlrs(t *testing.T) {
	prepTest()

	gid, _ := CreateGame(utils.DefGame)

	game := utils.DefGame
	game.GID = gid
	game.CurPlrs = 1

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(CurP) = "%v", want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(CurP) = %q, "%v", want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeCurHigh(t *testing.T) {
	prepTest()

	gid, _ := CreateGame(utils.DefGame)

	game := utils.DefGame
	game.GID = gid
	game.CurPlrs = 99

	err := UpdateGame(game)
	if err == nil {
		t.Fatalf(`UpdateGame(CurPHigh) = "%v", want ERROR`, err)
	}

	game = utils.DefGame
	game.GID = gid

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(CurPHigh) = %q, "%v", want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeAdditional(t *testing.T) {
	prepTest()

	gid, _ := CreateGame(utils.DefGame)

	game := utils.DefGame
	game.GID = gid
	game.AdditionalRules = `{"rule1": "default"}`

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(AddRules) = "%v", want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(AddRules) = %q, "%v", want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameChangeAll(t *testing.T) {
	prepTest()

	gid, _ := CreateGame(utils.DefGame)

	var game = utils.Game{
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
		t.Fatalf(`UpdateGame(All) = "%v", want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(All) = %q, "%v", want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameInvalidGID(t *testing.T) {
	prepTest()

	gid, _ := CreateGame(utils.DefGame)

	game := utils.DefGame
	game.GID = uuid.NewString()
	game.Name = "Game Test"

	err := UpdateGame(game)
	if err == nil {
		t.Fatalf(`UpdateGame(InvalidGID) = "%v", want ERROR`, err)
	}

	game = utils.DefGame
	game.GID = gid

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(InvalidGID) = %q, "%v", want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameEmptyGID(t *testing.T) {
	prepTest()

	gid, _ := CreateGame(utils.DefGame)

	game := utils.DefGame
	game.GID = ""
	game.Name = "Game Test"

	err := UpdateGame(game)
	if err == nil {
		t.Fatalf(`UpdateGame(EmptyGID) = "%v", want ERROR`, err)
	}

	game = utils.DefGame
	game.GID = gid

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(EmptyGID) = %q, "%v", want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameEmptyDB(t *testing.T) {
	prepTest()

	game := utils.DefGame
	game.GID = uuid.NewString()
	game.Name = "Game Test"

	err := UpdateGame(game)
	if err == nil {
		t.Fatalf(`UpdateGame(EmptyDB) = "%v", want ERROR`, err)
	}

	games, err := GetGames()
	if len(games) > 0 || err != nil {
		t.Fatalf(`UpdateGame(EmptyDB) - GET = %q, "%v", want "[]", "nil"`, games, err)
	}
}

func TestUpdateGameNoChange(t *testing.T) {
	prepTest()

	gid, _ := CreateGame(utils.DefGame)

	game := utils.DefGame
	game.GID = gid

	err := UpdateGame(game)
	if err != nil {
		t.Fatalf(`UpdateGame(NoChange) = "%v", want "nil"`, err)
	}

	games, err := GetGames()
	if game != games[0] || err != nil {
		t.Fatalf(`UpdateGame(NoChange) = %q, "%v", want %q, "nil"`, games, err, game)
	}
}

func TestUpdateGameInjection(t *testing.T) {
	prepTest()
	// TODO
}

func TestChangePlayerCountValidIncrease(t *testing.T) {
	prepTest()

	game := utils.DefGame
	game.CurPlrs = 2

	gid, _ := CreateGame(game)
	game.GID = gid

	err := ChangePlayerCount(game, 1)
	if err != nil {
		t.Fatalf(`ChangePlayerCount(ValidIncrease) = "%v", want "nil"`, err)
	}

	want := game.CurPlrs + 1

	game, err = GetGame(gid)
	if game.CurPlrs != want || err != nil {
		t.Fatalf(`ChangePlayerCount(ValidIncrease) = %d, "%v", want %d, "nil"`, game.CurPlrs, err, want)
	}
}

func TestChangePlayerCountValidDecrease(t *testing.T) {
	prepTest()

	game := utils.DefGame
	game.CurPlrs = 2

	gid, _ := CreateGame(utils.DefGame)

	game.GID = gid

	err := ChangePlayerCount(game, -1)
	if err != nil {
		t.Fatalf(`ChangePlayerCount(ValidDecrease) = "%v", want "nil"`, err)
	}

	want := game.CurPlrs - 1

	game, err = GetGame(gid)
	if game.CurPlrs != want || err != nil {
		t.Fatalf(`ChangePlayerCount(ValidDecrease) = %d, "%v", want %d, "nil"`, game.CurPlrs, err, want)
	}
}

func TestChangePlayerCountDoubleIncrease(t *testing.T) {
	prepTest()

	game := utils.DefGame
	game.CurPlrs = 1

	gid, _ := CreateGame(utils.DefGame)
	game.GID = gid

	want := game.CurPlrs + 2

	_ = ChangePlayerCount(game, 1)

	game, _ = GetGame(gid)

	err := ChangePlayerCount(game, 1)
	if err != nil {
		t.Fatalf(`ChangePlayerCount(DoubleIncrease) = "%v", want "nil"`, err)
	}

	game, err = GetGame(gid)
	if game.CurPlrs != want || err != nil {
		t.Fatalf(`ChangePlayerCount(DoubleIncrease) = %d, "%v", want %d, "nil"`, game.CurPlrs, err, want)
	}
}

func TestChangePlayerCountDoubleDecrease(t *testing.T) {
	prepTest()

	game := utils.DefGame
	game.CurPlrs = 3

	gid, _ := CreateGame(utils.DefGame)
	game.GID = gid

	want := game.CurPlrs - 2

	_ = ChangePlayerCount(game, -1)

	game, _ = GetGame(gid)

	err := ChangePlayerCount(game, -1)
	if err != nil {
		t.Fatalf(`ChangePlayerCount(DoubleDecrease) = "%v", want "nil"`, err)
	}

	game, err = GetGame(gid)
	if game.CurPlrs != want || err != nil {
		t.Fatalf(`ChangePlayerCount(DoubleDecrease) = %d, "%v", want %d, "nil"`, game.CurPlrs, err, want)
	}
}
func TestChangePlayerCountIncreaseFullGame(t *testing.T) {
	prepTest()

	game := utils.DefGame
	game.CurPlrs = game.MaxPlrs

	gid, _ := CreateGame(game)

	game.GID = gid

	err := ChangePlayerCount(game, 1)

	if err == nil {
		t.Fatalf(`ChangePlayerCount(FullGameIncrease) = "%v", want ERROR`, err)
	}

	want := game.CurPlrs

	game, err = GetGame(gid)
	if game.CurPlrs != want || err != nil {
		t.Fatalf(`ChangePlayerCount(FullGameIncrease) = %d, "%v", want %d, "nil"`, game.CurPlrs, err, want)
	}
}

func TestChangePlayerCountDecreaseToZero(t *testing.T) {
	prepTest()

	game := utils.DefGame
	game.CurPlrs = 1

	gid, _ := CreateGame(game)

	game.GID = gid

	err := ChangePlayerCount(game, -1)
	if err != nil {
		t.Fatalf(`ChangePlayerCount(DecreaseToZero) = "%v", want "nil"`, err)
	}

	game, err = GetGame(gid)
	if err == nil {
		t.Fatalf(`ChangePlayerCount(DecreaseToZero) = %q, "%v", want "nil", ERROR`, game, err)
	}
}

func TestChangePlayerCountNoChange(t *testing.T) {
	prepTest()

	game := utils.DefGame
	game.CurPlrs = 1

	gid, _ := CreateGame(game)

	game.GID = gid

	err := ChangePlayerCount(game, 0)
	if err != nil {
		t.Fatalf(`ChangePlayerCount(NoChange) = "%v", want "nil"`, err)
	}

	want := game.CurPlrs

	game, err = GetGame(gid)
	if game.CurPlrs != want || err != nil {
		t.Fatalf(`ChangePlayerCount(NoChange) %q, "%v", want %q, "nil"`, game.CurPlrs, err, want)
	}
}

// Test Game Deletion
func TestDeleteGameValid(t *testing.T) {
	prepTest()
	gid, _ := CreateGame(utils.DefGame)
	err := DeleteGame(gid)
	if err != nil {
		t.Fatalf(`DeleteGame(Valid) = "%v", want "nil"`, err)
	}

	games, err := GetGames()
	if len(games) > 0 || err != nil {
		t.Fatalf(`DeleteGame(Valid) = %q, "%v", want 0, "nil"`, len(games), err)
	}
}

func TestDeleteGameMulti(t *testing.T) {
	prepTest()
	_, _ = CreateGame(utils.DefGame)
	gid, _ := CreateGame(utils.DefGame)
	err := DeleteGame(gid)
	if err != nil {
		t.Fatalf(`DeleteGame(Multi) = "%v", want  "nil"`, err)
	}

	games, err := GetGames()
	if len(games) > 1 || err != nil {
		t.Fatalf(`DeletGame(Multi) = %q, "%v", want 1, "nil"`, len(games), err)
	}
}

func TestDeleteGameInvalidGID(t *testing.T) {
	prepTest()
	_, _ = CreateGame(utils.DefGame)

	err := DeleteGame(uuid.NewString())
	if err == nil {
		t.Fatalf(`DeleteGame(InvalidGID) = "%v", want ERROR`, err)
	}

	games, err := GetGames()
	if len(games) < 1 || err != nil {
		t.Fatalf(`DeletGame(InvalidGID) = %q, "%v", want 0, "nil"`, len(games), err)
	}
}

func TestDeletGameEmptyGID(t *testing.T) {
	prepTest()
	_, _ = CreateGame(utils.DefGame)

	err := DeleteGame("")
	if err == nil {
		t.Fatalf(`DeleteGame(EmptyGID) = "%v", want ERROR`, err)
	}

	games, err := GetGames()
	if len(games) < 1 || err != nil {
		t.Fatalf(`DeletGame(EmptyGID) = %q, "%v", want 0, "nil"`, len(games), err)
	}
}

func TestDeleteGameEmptyDB(t *testing.T) {
	prepTest()

	err := DeleteGame(uuid.NewString())
	if err == nil {
		t.Fatalf(`DeleteGame(EmptyDB) = "%v", want ERROR`, err)
	}
}

func TestDeleteGameRepeat(t *testing.T) {
	prepTest()
	gid, _ := CreateGame(utils.DefGame)
	_ = DeleteGame(gid)

	err := DeleteGame(gid)
	if err == nil {
		t.Fatalf(`DeleteGame(Repeat) = "%v", want ERROR`, err)
	}
}

func TestDeleteGameInjection(t *testing.T) {
	prepTest()
	//TODO

}
