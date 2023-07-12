package consequences

import (
	"Engee-Server/utils"
	"encoding/json"
	"testing"
	"time"

	db "Engee-Server/database"
	g "Engee-Server/game"
	u "Engee-Server/user"

	"github.com/google/go-cmp/cmp"
)

func prepareConGame(t *testing.T, testName string, userCount int) (string, string, []string) {
	db.InitDB()

	gid, err := g.CreateGame(defGame)
	if err != nil {
		t.Fatalf(`%v = failed to prepare conGame (creating game): %v`, testName, err)
	}

	plr := defPlr
	plr.GID = gid

	lid, _ := u.CreateUser(defPlr)
	plr.UID = lid
	err = u.UpdateUser(plr)
	if err != nil {
		t.Fatalf(`%v = failed to prepare conGame (updating leader): %v`, testName, err)
	}

	var users []string

	for i := 0; i < userCount; i++ {
		uid, _ := u.CreateUser(defPlr)
		plr.UID = uid
		err = u.UpdateUser(plr)
		if err != nil {
			t.Fatalf(`%v = failed to prepare conGame (updating user): %v`, testName, err)
		}

		users = append(users, uid)
	}

	return gid, lid, users
}

func initConGame(t *testing.T, testName string, userCount int) (string, string, []string) {
	gid, lid, users := prepareConGame(t, testName, userCount)

	initMsg := utils.GameMsg{
		Type:    "Init",
		UID:     lid,
		GID:     gid,
		Content: string(ts),
	}

	_, err := Handle(initMsg)
	if err != nil {
		t.Fatalf(`%v = "%v", want "nil"`, testName, err)
	}

	return gid, lid, users
}

func startConGame(t *testing.T, testName string, userCount int) (string, string, []string) {
	gid, lid, users := initConGame(t, testName, userCount)

	startMsg := utils.GameMsg{
		Type: "Start",
		UID:  lid,
		GID:  gid,
	}

	_, err := Handle(startMsg)
	if err != nil {
		t.Fatalf(`%v = "%v", want "nil"`, testName, err)
	}

	return gid, lid, users
}

func createWant(state string, users []string, timer int) ConVars {

	stories := map[string][]string{}
	for _, user := range users {
		stories[user] = []string{}
	}

	return ConVars{
		State:    state,
		Settings: testSettings,
		Timer:    timer,
		Stories:  stories,
	}
}

func checkTimeout(t *testing.T, testName string, timer int, gid string, want bool) {

	//TODO add toggle for timeouts
	if false {
		return
	}

	time.Sleep(time.Duration(timer) * time.Second)

	cVars, err := GetConState(gid)
	if err != nil {
		t.Fatalf(`%v = Could not get state: %v`, testName, err)
	}

	time := cVars.Timer
	if want {
		if time > 0 {
			t.Fatalf(`%v = "%d", want "0"`, testName, time)
		}
	} else {
		if time <= 0 {
			t.Fatalf(`%v = "%d", want ">0"`, testName, time)
		}
	}
}

func TestInitValid(t *testing.T) {
	lid, gid, _ := prepareConGame(t, "TestInit(Valid)", 3)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  lid,
		GID:  gid,
	}

	msg, err := Handle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestInit(Valid) = %q, "%v", want "ACK", "nil"`, msg.Type, err)
	}

	want := createWant("Lobby", []string{}, testSettings.Timer1)

	cVars, err := GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestInit(Valid) = %q, "%v", want %q, "nil"`, cVars, nil, want)
	}
}

func TestInitZeroRounds(t *testing.T) {
	lid, gid, _ := prepareConGame(t, "TestInit(ZeroRounds)", 3)

	rules := testSettings
	rules.Rounds = 0
	rs, _ := json.Marshal(rules)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  lid,
		GID:  gid,
	}

	msg, err := Handle(gMsg)

	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestInit(ZeroRounds) = %q, "%v", want "ACK", "nil"`, msg.Type, err)
	}
}

func TestInitNegativeRounds(t *testing.T) {
	lid, gid, _ := prepareConGame(t, "TestInit(NegativeRounds)", 3)

	rules := testSettings
	rules.Rounds = -1
	rs, _ := json.Marshal(rules)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  lid,
		GID:  gid,
	}

	msg, err := Handle(gMsg)

	if msg.Type != "Error" || err == nil {
		t.Fatalf(`TestInit(NegativeRounds) = %q, "%v", want "Error", ERROR`, msg.Type, err)
	}
}

func TestInitZeroShuffle(t *testing.T) {
	lid, gid, _ := prepareConGame(t, "TestInit(ZeroShuffle)", 3)

	rules := testSettings
	rules.Shuffle = 0
	rs, _ := json.Marshal(rules)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  lid,
		GID:  gid,
	}

	msg, err := Handle(gMsg)

	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestInit(ZeroShuffle) = %q, "%v", want "ACK", "nil"`, msg.Type, err)
	}
}

func TestInitInvalidShuffle(t *testing.T) {
	lid, gid, _ := prepareConGame(t, "TestInit(InvalidShuffle)", 3)

	rules := testSettings
	rules.Shuffle = 99
	rs, _ := json.Marshal(rules)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  lid,
		GID:  gid,
	}

	msg, err := Handle(gMsg)

	if msg.Type != "Error" || err == nil {
		t.Fatalf(`TestInit(InvalidShuffle) = %q, "%v", want "Error", Error`, msg.Type, err)
	}
}

func TestInitZeroTimer1(t *testing.T) {
	lid, gid, _ := prepareConGame(t, "TestInit(ZeroTimer)", 3)

	rules := testSettings
	rules.Timer1 = 0
	rs, _ := json.Marshal(rules)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  lid,
		GID:  gid,
	}

	msg, err := Handle(gMsg)

	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestInit(ZeroTimer) = %q, "%v", want "ACK", "nil"`, msg.Type, err)
	}
}

func TestNegativeTimer(t *testing.T) {
	lid, gid, _ := prepareConGame(t, "TestInit(ZeroTimer)", 3)

	rules := testSettings
	rules.Timer1 = -1
	rs, _ := json.Marshal(rules)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  lid,
		GID:  gid,
	}

	msg, err := Handle(gMsg)

	if msg.Type != "Error" || err == nil {
		t.Fatalf(`TestInit(ZeroTimer) = %q, "%v", want "Error", ERROR`, msg.Type, err)
	}
}

func TestInitZeroTimer2(t *testing.T) {
	lid, gid, _ := prepareConGame(t, "TestInit(ZeroTimer)", 3)

	rules := testSettings
	rules.Timer2 = 0
	rs, _ := json.Marshal(rules)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  lid,
		GID:  gid,
	}

	msg, err := Handle(gMsg)

	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestInit(ZeroTimer) = %q, "%v", want "ACK", "nil"`, msg.Type, err)
	}
}

func TestInitEmptyPrompts(t *testing.T) {
	lid, gid, _ := prepareConGame(t, "TestInit(EmptyPrompts)", 3)

	rules := testSettings
	rules.Prompts = []string{}
	rs, _ := json.Marshal(rules)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  lid,
		GID:  gid,
	}

	msg, err := Handle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestInit(EmptyPrompts) = %q, "%v", want "ACK", "nil"`, msg, err)
	}
}

func TestInitSinglePrompt(t *testing.T) {
	lid, gid, _ := prepareConGame(t, "TestInit(SinglePrompt)", 3)

	rules := testSettings
	rules.Prompts = []string{"Single Prompt"}
	rs, _ := json.Marshal(rules)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  lid,
		GID:  gid,
	}

	msg, err := Handle(gMsg)
	if msg.Type != "Error" || err == nil {
		t.Fatalf(`TestInit(SinglePrompts) = %q, "%v", want "Error", ERROR`, msg, err)
	}
}

func TestInitInvalidStruct(t *testing.T) {
	lid, gid, _ := prepareConGame(t, "TestInit(InvalidStruct)", 3)

	rs, _ := json.Marshal(defPlr)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  lid,
		GID:  gid,
	}

	msg, err := Handle(gMsg)
	if msg.Type != "Error" || err == nil {
		t.Fatalf(`TestInit(InvalidStruct) = %q, "%v", want "Error", ERROR`, msg, err)
	}
}

func TestInitEmpty(t *testing.T) {
	lid, gid, _ := prepareConGame(t, "TestInit(Invalid)", 3)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = ""
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  lid,
		GID:  gid,
	}

	msg, err := Handle(gMsg)
	if msg.Type != "Error" || err == nil {
		t.Fatalf(`TestInit(Valid) = %q, "%v", want "Error", ERROR`, msg, err)
	}
}

func TestStart(t *testing.T) {
	gid, lid, users := initConGame(t, "TestStart()", 3)

	gMsg := utils.GameMsg{
		Type: "Start",
		UID:  lid,
		GID:  gid,
	}

	msg, err := Handle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestStart() = %q, "%v", want "ACK", "err"`, msg.Type, err)
	}

	timer := testSettings.Timer1
	want := createWant("Prompts", users, timer)

	cVars, err := GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestPause() = %q, "%v", want %q, "nil"`, cVars, nil, want)
	}

}

func TestPauseTimer(t *testing.T) {
	gid, lid, users := startConGame(t, "TestPause()", 3)

	gMsg := utils.GameMsg{
		Type: "Pause",
		UID:  lid,
		GID:  gid,
	}

	msg, err := Handle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestPause() = %q, "%v", want "ACK", "nil"`, msg.Type, err)
	}

	want := createWant("Pause", users, testSettings.Timer1)

	cVars, err := GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestPause() = %q, "%v", want %q, "nil"`, cVars, nil, want)
	}

	checkTimeout(t, "TestPause()", testSettings.Timer1, gid, false)
}

func TestUnpauseTimer(t *testing.T) {
	gid, lid, users := startConGame(t, "TestUnpause()", 3)

	gMsg := utils.GameMsg{
		Type: "Pause",
		UID:  lid,
		GID:  gid,
	}

	_, _ = Handle(gMsg)

	msg, err := Handle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestUnpause() = %q, "%v", want "ACK", "nil"`, msg.Type, err)
	}

	want := createWant("Prompts", users, testSettings.Timer1)

	cVars, err := GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestUnpause() = %q, "%v", want %q, "nil"`, cVars, nil, want)
	}

	checkTimeout(t, "TestPause()", testSettings.Timer1, gid, true)
}

func TestEnd(t *testing.T) {
	gid, lid, _ := startConGame(t, "TestEnd()", 3)

	gMsg := utils.GameMsg{
		Type: "End",
		UID:  lid,
		GID:  gid,
	}

	msg, err := Handle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestEnd() = %q, "%v", want "ACK", "nil"`, msg.Type, err)
	}

	want := ConVars{}

	cVars, err := GetConState(gid)
	if !cmp.Equal(cVars, want) || err == nil {
		t.Fatalf(`TestEnd() = %q, "%v", want %q, ERROR`, cVars, nil, want)
	}
}

func TestReset(t *testing.T) {
	gid, lid, _ := startConGame(t, "TestReset()", 3)

	gMsg := utils.GameMsg{
		Type: "Reset",
		UID:  lid,
		GID:  gid,
	}

	msg, err := Handle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestReset() = %q, "%v", want "ACK", "nil"`, msg.Type, err)
	}

	want := createWant("Lobby", []string{}, testSettings.Timer1)

	cVars, err := GetConState(gid)
	if !cmp.Equal(cVars, want) || err == nil {
		t.Fatalf(`TestReset() = %q, "%v", want %q, ERROR`, cVars, nil, want)
	}

}

func TestRules(t *testing.T) {
	gid, lid, _ := startConGame(t, "TestRules()", 3)
	var ruleString string

	rules := testSettings
	rules.Rounds = 2

	gMsg := utils.GameMsg{
		Type:    "Rules",
		UID:     lid,
		GID:     gid,
		Content: ruleString,
	}

	msg, err := Handle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestRules() = %q, "%v", want "ACK", "nil"`, msg.Type, err)
	}

	want := createWant("Lobby", []string{}, testSettings.Timer1)

	cVars, err := GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestRules() = %q, "%v", want %q, "nil"`, cVars, nil, want)
	}
}

func TestRemove(t *testing.T) {
	gid, lid, users := startConGame(t, "TestRemove()", 3)

	gMsg := utils.GameMsg{
		Type:    "Remove",
		UID:     lid,
		GID:     gid,
		Content: users[0],
	}

	msg, err := Handle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestRemove() = %q, "%v", want "ACK", "nil"`, msg.Type, err)
	}

	users[0] = users[len(users)-1]

	want := createWant("Prompts", users, testSettings.Timer1)

	cVars, err := GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestRemove() = %q, "%v", want %q, "nil"`, cVars, nil, want)
	}
}

// Don't need to test remove too few, handled by leave below

func TestStatusConSpec(t *testing.T) {
	gid, lid, _ := startConGame(t, "TestStatus(ConSpec)", 3)

	status := "Replied"

	gMsg := utils.GameMsg{
		Type:    "Status",
		UID:     lid,
		GID:     gid,
		Content: status,
	}

	msg, err := Handle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestStatus(ConSpec) = %q, "%v", want "ACK", "nil"`, msg.Type, err)
	}

	plr, err := u.GetUser(lid)
	if plr.Status != status || err != nil {
		t.Fatalf(`TestStatus(ConSpec) = %q, "%v", want %q, "nil"`, plr.Status, err, status)
	}
}

func TestStatusPhaseChange(t *testing.T) {
	gid, lid, users := startConGame(t, "TestStatus(PhaseChange)", 3)

	gMsg := utils.GameMsg{
		Type:    "Status",
		UID:     lid,
		GID:     gid,
		Content: "Replied",
	}

	changeState("Stories", gid)

	for i, user := range users {
		gMsg.UID = user
		_, err := Handle(gMsg)
		if err != nil {
			t.Fatalf(`TestStatus(PhaseChange) [%d] = "%v", want "nil"`, i, err)
		}
	}

	want := createWant("Stories", users, testSettings.Timer2)
	gMsg.UID = lid

	msg, err := Handle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestStatus(PhaseChange) = %q, "%v", want "ACK", "nil"`, msg.Type, err)
	}

	cVars, err := GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestStatus(PhaseChange) = %q, "%v", want "ACK", "nil"`, cVars, err)
	}

	plr, err := u.GetUser(lid)
	if plr.Status != "Reading" || err != nil {
		t.Fatalf(`TestStatus(ConSpec) = %q, "%v", want "Reading, "nil"`, plr.Status, err)
	}
}

func TestStatusInvalid(t *testing.T) {
	gid, lid, _ := startConGame(t, "TestStatus(Invalid)", 3)

	gMsg := utils.GameMsg{
		Type:    "Status",
		UID:     lid,
		GID:     gid,
		Content: "Invalid",
	}

	msg, err := Handle(gMsg)
	if msg.Type != "Error" || err == nil {
		t.Fatalf(`TestStatus(Invalid) = %q, "%v", want "Error", ERROR`, msg.Type, err)
	}

	plr, err := u.GetUser(lid)
	if plr.Status != "Ready" || err != nil {
		t.Fatalf(`TestStatus(Invalid) = %q, "%v", want "Ready", "nil"`, plr.Status, err)
	}
}

func TestLeaveValid(t *testing.T) {
	gid, lid, users := startConGame(t, "TestLeave(Valid)", 3)

	gMsg := utils.GameMsg{
		Type: "Leave",
		UID:  lid,
		GID:  gid,
	}

	want := createWant("Prompts", users, testSettings.Timer1)
	delete(want.Stories, users[0])

	msg, err := Handle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestLeave(Valid) = %q, "%v", want "ACK", "nil"`, msg.Type, err)
	}

	cVars, err := GetConState(gid)
	if len(cVars.Stories) != len(want.Stories) || err != nil {
		t.Fatalf(`TestLeave(Valid) = %q, "%v", want %q, "nil"`, cVars.Stories, err, want.Stories)
	}
}

func TestLeaveTooFew(t *testing.T) {
	gid, _, users := startConGame(t, "TestLeave(TooFewUsers)", 2)

	gMsg := utils.GameMsg{
		Type: "Leave",
		UID:  users[0],
		GID:  gid,
	}

	msg, err := Handle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestLeave(TooFew) = %q, "%v", want "ACK", "nil"`, msg.Type, err)
	}

	want := ConVars{}

	cVars, err := GetConState(gid)
	if !cmp.Equal(cVars, want) || err == nil {
		t.Fatalf(`TestLeave(TooFew) = %q, "%v", want %q, ERROR`, cVars, err, want)
	}

}

func TestReplyValid(t *testing.T) {
	gid, lid, users := startConGame(t, "TestReply(Valid)", 2)

	story, _ := json.Marshal(defStory)

	gMsg := utils.GameMsg{
		Type:    "Story",
		UID:     lid,
		GID:     gid,
		Content: string(story),
	}

	msg, err := Handle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestReply(Valid) = %q "%v", want "ACK", "nil"`, msg, err)
	}

	want := createWant("Prompts", users, testSettings.Timer1)
	want.Stories[lid] = defStory

	cVars, err := GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestReply(Valid) = %q, "%v", want %q, "nil"`, cVars, err, want)
	}
}

func TestReplyPhaseChange(t *testing.T) {
	gid, lid, users := startConGame(t, "TestReply(Valid)", 2)

	story, _ := json.Marshal(defStory)

	gMsg := utils.GameMsg{
		Type:    "Story",
		UID:     lid,
		GID:     gid,
		Content: string(story),
	}

	msg, err := Handle(gMsg)
	if msg.Type != "ACK" || err != nil {
		t.Fatalf(`TestReply(Valid) = %q "%v", want "ACK", "nil"`, msg, err)
	}

	want := createWant("Stories", users, testSettings.Timer1)
	want.Stories[lid] = defStory

	cVars, err := GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestReply(Valid) = %q, "%v", want %q, "nil"`, cVars, err, want)
	}

}

func TestReplyShort(t *testing.T) {
	gid, lid, users := startConGame(t, "TestReply(Short)", 2)

	wrongStory := defStory[0:3]
	story, _ := json.Marshal(wrongStory)

	gMsg := utils.GameMsg{
		Type:    "Story",
		UID:     lid,
		GID:     gid,
		Content: string(story),
	}

	msg, err := Handle(gMsg)
	if msg.Type != "Error" || err == nil {
		t.Fatalf(`TestReply(Short) = %q "%v", want "Error", ERROR`, msg, err)
	}

	want := createWant("Prompts", users, testSettings.Timer1)

	cVars, err := GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestReply(Short) = %q, "%v", want %q, "nil"`, cVars, err, want)
	}
}

func TestReplyLong(t *testing.T) {
	gid, lid, users := startConGame(t, "TestReply(Long)", 2)

	wrongStory := append(defStory, "Invalid")
	story, _ := json.Marshal(wrongStory)

	gMsg := utils.GameMsg{
		Type:    "Story",
		UID:     lid,
		GID:     gid,
		Content: string(story),
	}

	msg, err := Handle(gMsg)
	if msg.Type != "Error" || err == nil {
		t.Fatalf(`TestReply(Long) = %q "%v", want "Error", ERROR`, msg, err)
	}

	want := createWant("Prompts", users, testSettings.Timer1)

	cVars, err := GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestReply(Long) = %q, "%v", want %q, "nil"`, cVars, err, want)
	}
}

func TestReplyEmptyLine(t *testing.T) {
	gid, lid, users := startConGame(t, "TestReply(EmptyLine)", 2)

	wrongStory := defStory
	wrongStory[2] = ""
	story, _ := json.Marshal(wrongStory)

	gMsg := utils.GameMsg{
		Type:    "Story",
		UID:     lid,
		GID:     gid,
		Content: string(story),
	}

	msg, err := Handle(gMsg)
	if msg.Type != "Error" || err == nil {
		t.Fatalf(`TestReply(EmptyLine) = %q "%v", want "Error", ERROR`, msg, err)
	}

	want := createWant("Prompts", users, testSettings.Timer1)

	cVars, err := GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestReply(EmptyLine) = %q, "%v", want %q, "nil"`, cVars, err, want)
	}
}
