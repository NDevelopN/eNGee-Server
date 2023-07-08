package gamespace

import (
	"log"
	"testing"

	db "Engee-Server/database"
	g "Engee-Server/game"
	u "Engee-Server/user"
	utils "Engee-Server/utils"

	"github.com/google/uuid"
)

func PrepareLeaderTest(t *testing.T, pCount int) (string, string, []string, error) {
	db.InitDB()

	lid, _ := u.CreateUser(utils.DefUser)
	game := utils.DefGame
	game.Leader = lid

	gid, _ := g.CreateGame(game)

	leader := utils.DefLeader
	leader.UID = lid
	leader.GID = gid

	u.UpdateUser(leader)

	_ = ChangeStatus(lid, gid, "Ready")

	user := utils.DefUser
	user.GID = gid

	var users []string

	for i := 0; i < pCount; i++ {
		p, _ := u.CreateUser(utils.DefUser)
		user.UID = p
		err := u.UpdateUser(user)
		if err != nil {
			log.Fatalf("Failed to update user: %v", err)
		}

		users = append(users, p)
		_ = ChangeStatus(p, gid, "Ready")
	}

	return gid, lid, users, nil
}

func CheckFinalStatus(t *testing.T, gid string, testName string, want string) {
	game, err := g.GetGame(gid)
	if err != nil {
		t.Fatalf(`%v(Getting Game) = %q, want "nil"`, testName, err)
	}

	if game.Status != want {
		t.Fatalf(`%v = Game status: %q, want %q`, testName, game.Status, want)
	}
}

func TestPauseValid(t *testing.T) {
	gid, lid, _, err := PrepareLeaderTest(t, 0)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestPause(Valid) = Failed to prepare for test: %q`, err)
	}

	err = Pause(gid, lid)
	if err != nil {
		t.Fatalf(`TestPause(Valid) = %q, want "nil"`, err)
	}

	CheckFinalStatus(t, gid, "TestPause(Valid)", "Pause")
}

func TestUnPauseValid(t *testing.T) {
	gid, lid, _, err := PrepareLeaderTest(t, 0)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestUnPause(Valid) = Failed to prepare for test: %q`, err)
	}

	_ = Pause(gid, lid)

	err = Pause(gid, lid)
	if err != nil {
		t.Fatalf(`TestUnPause(Valid) = %q, want "nil"`, err)
	}

	CheckFinalStatus(t, gid, "TestUnPause(Valid)", "Lobby")
}
func TestPauseMidGame(t *testing.T) {
	gid, lid, _, err := PrepareLeaderTest(t, 2)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestPause(MidGame) = Failed to prepare for test: %q`, err)
	}

	_ = Start(gid, lid)

	err = Pause(gid, lid)
	if err != nil {
		t.Fatalf(`TestPause(MidGame) = %q, want "nil"`, err)
	}

	CheckFinalStatus(t, gid, "TestPause(MidGame)", "Pause")
}

func TestUnPauseMidGame(t *testing.T) {
	gid, lid, _, err := PrepareLeaderTest(t, 2)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestUnPause(MidGame) = Failed to prepare for test: %q`, err)
	}

	_ = Start(gid, lid)
	_ = Pause(gid, lid)

	err = Pause(gid, lid)
	if err != nil {
		t.Fatalf(`TestUnPause(MidGame) = %q, want "nil"`, err)
	}

	CheckFinalStatus(t, gid, "TestUnPause(MidGame)", "Play")
}
func TestPauseNotLeader(t *testing.T) {
	gid, lid, _, err := PrepareLeaderTest(t, 0)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestPause(NotLeader) = Failed to prepare for test: %q`, err)
	}

	err = Pause(gid, uuid.NewString())
	if err == nil {
		t.Fatalf(`TestPause(NotLeader) = %q, want ERROR`, err)
	}

	CheckFinalStatus(t, gid, "TestPause(NotLeader)", "Lobby")
}

func TestStartValid(t *testing.T) {
	gid, lid, _, err := PrepareLeaderTest(t, 2)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestStart(Valid) = Failed to prepare for test: %q`, err)
	}

	err = Start(gid, lid)
	if err != nil {
		t.Fatalf(`TestStart(Valid) = %q, want "nil"`, err)
	}

	CheckFinalStatus(t, gid, "TestStart(Valid)", "Play")
}

func TestStartNotLeader(t *testing.T) {
	gid, lid, _, err := PrepareLeaderTest(t, 2)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestStart(Not Leader) = Failed to prepare for test: %q`, err)
	}

	err = Start(gid, uuid.NewString())
	if err == nil {
		t.Fatalf(`TestStart(Not Leader) = %q, want ERROR`, err)
	}

	CheckFinalStatus(t, gid, "TestStart(NotLeader)", "Lobby")
}

func TestStartNotEnoughPlayers(t *testing.T) {
	gid, lid, _, err := PrepareLeaderTest(t, 1)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestStart(NotEnoughPlrs) = Failed  to prepare for test: %q`, err)
	}

	err = Start(gid, lid)
	if err == nil {
		t.Fatalf(`TestStart(NotEnoughPlayers) = %q, want ERROR`, err)
	}

	CheckFinalStatus(t, gid, "TestStart(NotEnoughPlrs)", "Lobby")
}

func TestStartNonEnoughReady(t *testing.T) {
	gid, lid, users, err := PrepareLeaderTest(t, 2)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestStart(NotEnoughReady) = Failed  to prepare for test: %q`, err)
	}

	for _, user := range users {
		_ = ChangeStatus(user, gid, "Not Ready")
	}

	err = Start(gid, lid)
	if err == nil {
		t.Fatalf(`TestStart(NotEnoughReady) = %q, want ERROR`, err)
	}

	CheckFinalStatus(t, gid, "TestStart(NotEnoughReady)", "Lobby")
}

func TestStartMidGame(t *testing.T) {
	gid, lid, _, err := PrepareLeaderTest(t, 2)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestStart(MidGame) = Failed to prepare for test: %q`, err)
	}

	_ = Start(gid, lid)

	err = Start(gid, lid)
	if err == nil {
		t.Fatalf(`TestStart(MidGame) = %q, want ERROR`, err)
	}

	CheckFinalStatus(t, gid, "TestStart(MidGame)", "Play")
}

func TestResetValid(t *testing.T) {
	gid, lid, _, err := PrepareLeaderTest(t, 0)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestReset(Valid) = Failed to prepare for test: %q`, err)
	}

	err = Reset(gid, lid)
	if err != nil {
		t.Fatalf(`TestReset(Valid) = %q, want "nil"`, err)
	}

	CheckFinalStatus(t, gid, "TestReset(Valid)", "Lobby")
}

func TestResetNotLeader(t *testing.T) {
	gid, lid, _, err := PrepareLeaderTest(t, 0)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestReset(NotLeader) = Failed to prepare for test: %q`, err)
	}

	err = Reset(gid, uuid.NewString())
	if err == nil {
		t.Fatalf(`TestReset(NotLeader) = %q, want ERROR`, err)
	}

	CheckFinalStatus(t, gid, "TestReset(NotLeader)", "Lobby")
}

func TestResetMidGame(t *testing.T) {
	gid, lid, _, err := PrepareLeaderTest(t, 2)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestReset(Valid-MidGame) = Failed to prepare for test: %q`, err)
	}

	_ = Start(gid, lid)

	err = Reset(gid, lid)
	if err != nil {
		t.Fatalf(`TestReset(Valid-MidGame) = %q, want "nil"`, err)
	}

	CheckFinalStatus(t, gid, "TestReset(MidGame)", "Lobby")
}

func TestEndValid(t *testing.T) {
	gid, lid, _, err := PrepareLeaderTest(t, 0)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestEnd(Valid) = Failed to prepare for test: %q`, err)
	}

	err = End(gid, lid)
	if err != nil {
		t.Fatalf(`TestEnd(Valid) = %q, want "nil"`, err)
	}

	game, err := g.GetGame(gid)
	if err == nil {
		t.Fatalf(`TestEnd(Valid) = %q, %q, want "", ERROR`, game, err)
	}
}

func TestEndNotLeader(t *testing.T) {
	gid, lid, _, err := PrepareLeaderTest(t, 0)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestEnd(Not Leader) = Failed to prepare for test: %q`, err)
	}

	err = End(gid, uuid.NewString())
	if err == nil {
		t.Fatalf(`TestEnd(Not Leader) = %q, want ERROR`, err)
	}

	game, err := g.GetGame(gid)
	if err != nil {
		t.Fatalf(`TestEnd(NotLeader) = %q, %q, want (Game), "nil"`, game, err)
	}
}

func TestEndMidGame(t *testing.T) {
	gid, lid, _, err := PrepareLeaderTest(t, 2)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestEnd(Valid-MidGame) = Failed to prepare for test: %q`, err)
	}

	_ = Start(gid, lid)

	err = End(gid, lid)
	if err != nil {
		t.Fatalf(`TestEnd(Valid-MidGame) = %q, want "nil"`, err)
	}

	game, err := g.GetGame(gid)
	if err == nil {
		t.Fatalf(`TestEnd(Valid) = %q, %q, want "", ERROR`, game, err)
	}
}

func TestRemoveValid(t *testing.T) {
	gid, lid, users, err := PrepareLeaderTest(t, 4)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestRemove(Valid) = Failed to prepare for test: %q`, err)
	}

	err = Remove(gid, lid, users[0])
	if err != nil {
		t.Fatalf(`TestRemove(Valid) = %q, want "nil"`, err)
	}

	CheckFinalStatus(t, gid, "TestRemove(Valid)", "Lobby")
}

func TestRemoveNotLeader(t *testing.T) {
	gid, lid, users, err := PrepareLeaderTest(t, 4)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestRemove(NotLeader) = Failed to prepare for test: %q`, err)
	}

	err = Remove(gid, users[0], users[0])
	if err == nil {
		t.Fatalf(`TestRemove(NotLeader) = %q, want ERROR`, err)
	}

	CheckFinalStatus(t, gid, "TestRemove(NotLeader)", "Lobby")
}

func TestRemoveInvalidPlayer(t *testing.T) {
	gid, lid, _, err := PrepareLeaderTest(t, 4)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestRemove(InvalidPlayer) = Failed to prepare for test: %q`, err)
	}

	err = Remove(gid, lid, uuid.NewString())
	if err == nil {
		t.Fatalf(`TestRemove(InvalidPlayer) = %q, want ERROR`, err)
	}

	CheckFinalStatus(t, gid, "TestRemove(InvalidPlayer)", "Lobby")
}

func TestRemoveSelf(t *testing.T) {
	gid, lid, _, err := PrepareLeaderTest(t, 4)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestRemove(Self) = Failed to prepare for test: %q`, err)
	}

	err = Remove(gid, lid, lid)
	if err == nil {
		t.Fatalf(`TestRemove(Self) = %q, want ERROR`, err)
	}

	CheckFinalStatus(t, gid, "TestRemove(Self)", "Lobby")
}

func TestRemoveMidGame(t *testing.T) {
	gid, lid, users, err := PrepareLeaderTest(t, 4)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestRemove(Valid-MidGame) = Failed to prepare for test: %q`, err)
	}

	_ = Start(gid, lid)

	err = Remove(gid, lid, users[0])
	if err != nil {
		t.Fatalf(`TestRemove(Valid-MidGame) = %q, want "nil"`, err)
	}

	game, _ := g.GetGame(gid)

	if game.Status != "Play" {
		t.Fatalf(`TestRemove(TooFewPlrs) = Game status: %q, want "Play"`, game.Status)
	}

	CheckFinalStatus(t, gid, "TestRemove(MidGame)", "Play")
}

func TestRemoveTooFewPlayers(t *testing.T) {
	gid, lid, users, err := PrepareLeaderTest(t, 2)
	if gid == "" || lid == "" || err != nil {
		t.Fatalf(`TestRemove(TooFewPlrs) = Failed to prepare for test: %q`, err)
	}

	_ = Start(gid, lid)

	err = Remove(gid, lid, users[0])
	if err != nil {
		t.Fatalf(`TestRemove(TooFewPlrs) = %q, want "nil"`, err)
	}

	CheckFinalStatus(t, gid, "TestRemove(MidGameFewPlrs)", "Lobby")
}
