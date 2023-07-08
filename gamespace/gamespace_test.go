package gamespace

import (
	db "Engee-Server/database"
	g "Engee-Server/game"
	u "Engee-Server/user"
	utils "Engee-Server/utils"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func PrepareGame(t *testing.T, testName string) (string, string) {
	db.InitDB()

	gid, err := g.CreateGame(utils.DefGame)
	if err != nil {
		t.Fatalf("%v = failed to prepare game (creating game): %v", testName, err)
	}

	plr := utils.DefUser
	plr.GID = gid
	plr.Status = "Ready"

	var uid string

	lid, _ := u.CreateUser(utils.DefUser)
	plr.UID = lid
	err = u.UpdateUser(plr)
	if err != nil {
		t.Fatalf("%v = failed to prepare game (updating leader): %v", testName, err)
	}

	for i := 0; i < 3; i++ {
		uid, _ = u.CreateUser(utils.DefUser)
		plr.UID = uid
		err = u.UpdateUser(plr)
		if err != nil {
			t.Fatalf("%v = failed to prepare game (updating user): %v", testName, err)
		}
	}

	return gid, lid
}

func PrepareInitTest(t *testing.T, testName string) utils.GameMsg {
	gid, uid := PrepareGame(t, testName)

	ts, err := json.Marshal(utils.TestSettings)
	if err != nil {
		t.Fatalf("%v = failed to prepare initialization test (marshalling settings): %v", testName, err)
	}

	gMsg := utils.GameMsg{
		Type:    "Init",
		UID:     uid,
		GID:     gid,
		Content: string(ts),
	}

	return gMsg
}

// Test Initialization
func TestInitValid(t *testing.T) {
	gMsg := PrepareInitTest(t, "TestInit(Valid)")
	msg, err := GamespaceHandle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestInit(Valid) = %q, "%v", want "ACK", "nil"`, msg.Type, err)
	}
}

func TestInitMultiGames(t *testing.T) {
	gMsg := PrepareInitTest(t, "TestInit(Multi)")
	_, _ = g.CreateGame(utils.DefGame)

	msg, err := GamespaceHandle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestInit(Multi) = %q, "%v", want "ACK", "nil"`, msg.Type, err)
	}
}

func TestInitInvalidSettings(t *testing.T) {
	gMsg := PrepareInitTest(t, "TestInit(InvalidSettings)")
	gMsg.Content = `"invalid": "content"`
	msg, err := GamespaceHandle(gMsg)
	if msg.Type != "Error" || err == nil {
		t.Fatalf(`TestInit(InvalidSettings) = %q, "%v", want "Error", ERROR`, msg.Type, err)
	}
}

// Test handling pre-checks
func TestHandleInvalidType(t *testing.T) {
	gMsg := PrepareInitTest(t, "TestHandle(InvalidType)")
	gMsg.Type = "Invalid"
	msg, err := GamespaceHandle(gMsg)
	if msg.Type != "Error" || err == nil {
		t.Fatalf(`TestHandle(InvalidType) = %q, "%v", want "Error", ERROR`, msg.Type, err)
	}
}

func TestHandleInvalidGID(t *testing.T) {
	gMsg := PrepareInitTest(t, "TestHandle(InvalidGID)")
	gMsg.GID = uuid.NewString()
	msg, err := GamespaceHandle(gMsg)
	if msg.Type != "Error" || err == nil {
		t.Fatalf(`TestHandle(InvalidGID) = %q, "%v", want "Error", ERROR`, msg.Type, err)
	}
}

func TestHandleInvalidUID(t *testing.T) {
	gMsg := PrepareInitTest(t, "TestHandle(InvalidUID)")
	gMsg.UID = uuid.NewString()
	msg, err := GamespaceHandle(gMsg)
	if msg.Type != "Error" || err == nil {
		t.Fatalf(`TestHandle(InvalidUID) = %q, "%v", want "Error", ERROR`, msg.Type, err)
	}
}

func TestHandleNotLeader(t *testing.T) {
	gMsg := PrepareInitTest(t, "TestHandle(NotLeader)")
	user := utils.DefUser
	user.GID = gMsg.GID
	gMsg.UID, _ = u.CreateUser(user)

	msg, err := GamespaceHandle(gMsg)
	if msg.Type != "Error" || err == nil {
		t.Fatalf(`TestHandle(NotLeader) = %q, "%v", want "Error", ERROR`, msg.Type, err)
	}
}

// Test Pause handling
func TestHandlePause(t *testing.T) {
	gid, uid := PrepareGame(t, "TestHandlePause(Valid)")

	gMsg := utils.GameMsg{
		Type: "Pause",
		GID:  gid,
		UID:  uid,
	}

	msg, err := GamespaceHandle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestHandlePause() = %q, "%v", want "ACK", "nil"`, msg.Type, err)
	}

	game, err := g.GetGame(gid)
	if game.Status != "Pause" || err != nil {
		t.Fatalf(`TestHandlePause() = %q, "%v", want "Play", "nil"`, game.Status, err)
	}
}

func TestHandleUnpause(t *testing.T) {
	gid, uid := PrepareGame(t, "TestHandleUnpause(Valid)")

	gMsg := utils.GameMsg{
		Type: "Pause",
		GID:  gid,
		UID:  uid,
	}

	_, _ = GamespaceHandle(gMsg)
	msg, err := GamespaceHandle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestHandleUnpause() = %q, "%v", want "ACK", "nil"`, msg.Type, err)
	}

	game, err := g.GetGame(gid)
	if game.Status != "Lobby" || err != nil {
		t.Fatalf(`TestHandleUnpause() = %q, "%v", want "Lobby", "nil"`, game.Status, err)
	}
}

// Test End handling
func TestHandleEnd(t *testing.T) {
	gid, uid := PrepareGame(t, "TestHandleEnd(Valid)")

	gMsg := utils.GameMsg{
		Type: "End",
		GID:  gid,
		UID:  uid,
	}

	msg, err := GamespaceHandle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestHandleEnd() = %q, "%v", want "ACK`, msg, err)
	}

	game, err := g.GetGame(gid)
	if err == nil {
		t.Fatalf(`TestHandleEnd() = %q, "%v", want nil, ERROR`, game, err)
	}
}

func TestHandleStart(t *testing.T) {
	gid, uid := PrepareGame(t, "TestHandleStart(Valid)")

	gMsg := utils.GameMsg{
		Type: "Start",
		GID:  gid,
		UID:  uid,
	}

	msg, err := GamespaceHandle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestHandleStart() = %q, "%v", want nil, ERROR`, msg, err)
	}

	game, err := g.GetGame(gid)
	if game.Status != "Play" || err != nil {
		t.Fatalf(`TestHandleStart() = %q, "%v", want "Play", "nil"`, game.Status, err)
	}
}

func TestHandleReset(t *testing.T) {
	gid, uid := PrepareGame(t, "TestHandleReset(Valid)")

	gMsg := utils.GameMsg{
		Type: "Start",
		GID:  gid,
		UID:  uid,
	}

	_, _ = GamespaceHandle(gMsg)

	gMsg.Type = "Reset"
	msg, err := GamespaceHandle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestHandleReset() = %q, "%v", want nil, ERROR`, msg, err)
	}

	game, err := g.GetGame(gid)
	if game.Status != "Lobby" || err != nil {
		t.Fatalf(`TestHandleReset() = %q, "%v", want "Lobby", "nil"`, game.Status, err)
	}
}

func TestHandleRules(t *testing.T) {
	gid, uid := PrepareGame(t, "TestHandleInvalidType")

	_, _ = g.GetGame(gid)

	ruleStruct := utils.DefGame
	ruleStruct.GID = gid
	ruleStruct.Leader = uid
	ruleStruct.MinPlrs = 2

	rules, _ := json.Marshal(ruleStruct)

	gMsg := utils.GameMsg{
		Type:    "Rules",
		GID:     gid,
		UID:     uid,
		Content: string(rules),
	}

	msg, err := GamespaceHandle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestHandleRules() = %q, "%v", want "ACK", "nil"`, msg, err)
	}

	game, err := g.GetGame(gid)
	if game != ruleStruct || err != nil {
		t.Fatalf(`TestHandleRules() = %q, "%v", want %q, "nil"`, game, err, ruleStruct)
	}
}

func TestHandleRemove(t *testing.T) {
	gid, uid := PrepareGame(t, "TestHandleInvalidType")

	target, _ := u.CreateUser(utils.DefUser)

	user := utils.DefUser
	user.UID = target
	user.GID = gid

	_ = u.UpdateUser(user)

	gMsg := utils.GameMsg{
		Type:    "Remove",
		GID:     gid,
		UID:     uid,
		Content: target,
	}

	msg, err := GamespaceHandle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestHandleRemove() = %q, "%v", want "ACK", "nil"`, msg, err)
	}
}

func TestHandleStatus(t *testing.T) {
	gid, uid := PrepareGame(t, "TestHandleInvalidType")

	gMsg := utils.GameMsg{
		Type:    "Status",
		GID:     gid,
		UID:     uid,
		Content: "Ready",
	}

	msg, err := GamespaceHandle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestHandleStatus() = %q, "%v", want "ACK", "nil"`, msg, err)
	}

	user, err := u.GetUser(uid)
	if user.Status != "Ready" || err != nil {
		t.Fatalf(`TestHandleStatus() = %q, "%v", want "Ready", "nil"`, msg, err)
	}
}

func TestHandleLeave(t *testing.T) {
	gid, _ := PrepareGame(t, "TestHandleInvalidType")

	user := utils.DefUser
	user.GID = gid

	uid, _ := u.CreateUser(utils.DefUser)

	user.UID = uid

	_ = u.UpdateUser(user)

	gMsg := utils.GameMsg{
		Type: "Leave",
		GID:  gid,
		UID:  uid,
	}

	msg, err := GamespaceHandle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestHandleLeave() = %q, "%v", want "ACK", "nil"`, msg, err)
	}

	user, err = u.GetUser(uid)
	if user.GID != "" || err != nil {
		t.Fatalf(`TestHandleLeave() = %q, "%v", want "", "nil"`, user.GID, err)
	}
}
