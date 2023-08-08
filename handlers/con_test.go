package handlers

import (
	"Engee-Server/utils"
	"encoding/json"
	"log"
	"testing"
	"time"

	c "Engee-Server/consequences"
	db "Engee-Server/database"
	g "Engee-Server/game"
	u "Engee-Server/user"

	"github.com/google/go-cmp/cmp"
)

const waitForTimeout = true

func prepareConGame(t *testing.T, testName string, userCount int) (string, []string) {
	db.InitDB()
	Init()

	utils.SETLOCALTEST(true)

	c.CVars = make(map[string]c.ConVars)

	lid, err := u.CreateUser(c.DefPlr)
	if err != nil {
		t.Fatalf(`%v = failed to prepare conGame (creating leader): %v`, testName, err)
	}

	game := c.DefGame
	game.Leader = lid

	gid, err := g.CreateGame(c.DefGame)
	if err != nil {
		t.Fatalf(`%v = failed to prepare conGame (creating game): %v`, testName, err)
	}

	var users = []string{lid}

	plr := c.DefPlr
	plr.UID = lid
	plr.GID = gid
	plr.Status = "New"

	err = u.UpdateUser(plr)
	if err != nil {
		t.Fatalf(`%v = failed to prepare conGame (updating leader): %v`, testName, err)
	}

	for i := 1; i < userCount; i++ {
		uid, err := u.CreateUser(c.DefPlr)
		if err != nil {
			t.Fatalf(`%v = failed to prepare conGame (creating user): %v`, testName, err)
		}

		plr.UID = uid
		err = u.UpdateUser(plr)
		if err != nil {
			t.Fatalf(`%v = failed to prepare conGame (updating user): %v`, testName, err)
		}

		users = append(users, uid)
	}

	log.Printf("Leader: %s", game.Leader)
	for i, k := range users {
		log.Printf("User %d: %s", i, k)
	}

	return gid, users
}

func initConGame(t *testing.T, testName string, userCount int) (string, []string) {
	gid, users := prepareConGame(t, testName, userCount)

	initMsg := utils.GameMsg{
		Type:    "Init",
		UID:     users[0],
		GID:     gid,
		Content: string(c.Ts),
	}

	cause, resp := c.Handle(initMsg)
	if cause != "" {
		t.Fatalf(`%v = %q - %q, want "" - ""`, testName, cause, resp)
	}

	return gid, users
}

func startConGame(t *testing.T, testName string, userCount int) (string, []string) {
	gid, users := initConGame(t, testName, userCount)

	startMsg := utils.GameMsg{
		Type: "Start",
		UID:  users[0],
		GID:  gid,
	}

	cause, resp := c.Handle(startMsg)
	if cause != "" {
		t.Fatalf(`%v = %q - %q, want "" - ""`, testName, cause, resp)
	}

	return gid, users
}

func createWant(state int, users []string, timer int) c.ConVars {

	stories := map[string][]string{}
	if len(users) > 0 {
		for _, user := range users {
			stories[user] = []string{}
		}
	}

	return c.ConVars{
		State:    state,
		Settings: c.TestSettings,
		Timer:    timer,
		Stories:  stories,
	}
}

func checkTimeout(t *testing.T, testName string, timer int, gid string, want bool) {
	if !waitForTimeout {
		return
	}

	time.Sleep(time.Duration(timer) * time.Second)

	cVars, err := c.GetConState(gid)
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
	gid, users := prepareConGame(t, "TestInit(Valid)", 3)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  users[0],
		GID:  gid,
	}

	cause, resp := c.Handle(gMsg)
	if cause != "" {
		t.Fatalf(`TestInit(Valid) = %q - %q, want "" - ""`, cause, resp)
	}

	want := createWant(c.LOBBY, nil, c.TestSettings.Timer1)

	cVars, err := c.GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestInit(Valid) = %v, "%v", want %v, "nil"`, cVars, err, want)
	}
}

func TestInitZeroRounds(t *testing.T) {
	gid, users := prepareConGame(t, "TestInit(ZeroRounds)", 3)

	rules := c.TestSettings
	rules.Rounds = 0
	rs, _ := json.Marshal(rules)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  users[0],
		GID:  gid,
	}

	cause, resp := c.Handle(gMsg)
	if cause != "" {
		t.Fatalf(`TestInit(ZeroRounds) = %q - %q, want "" - ""`, cause, resp)
	}
}

func TestInitNegativeRounds(t *testing.T) {
	gid, users := prepareConGame(t, "TestInit(NegativeRounds)", 3)

	rules := c.TestSettings
	rules.Rounds = -1
	rs, _ := json.Marshal(rules)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  users[0],
		GID:  gid,
	}

	want := "[Err Msg]"

	cause, resp := c.Handle(gMsg)
	if cause != "Error" {
		t.Fatalf(`TestInit(NegativeRounds) = %q - %q, want "Error" - %q`, cause, resp, want)
	}
}

func TestInitZeroShuffle(t *testing.T) {
	gid, users := prepareConGame(t, "TestInit(ZeroShuffle)", 3)

	rules := c.TestSettings
	rules.Shuffle = 0
	rs, _ := json.Marshal(rules)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  users[0],
		GID:  gid,
	}

	cause, resp := c.Handle(gMsg)
	if cause != "" {
		t.Fatalf(`TestInit(ZeroShuffle) = %q - %q, want "" - ""`, cause, resp)
	}
}

func TestInitInvalidShuffle(t *testing.T) {
	gid, users := prepareConGame(t, "TestInit(InvalidShuffle)", 3)

	rules := c.TestSettings
	rules.Shuffle = 99
	rs, _ := json.Marshal(rules)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  users[0],
		GID:  gid,
	}

	want := "[Err Msg]"

	cause, resp := c.Handle(gMsg)
	if cause != "Error" {
		t.Fatalf(`TestInit(InvalidShuffle) = %q - %q, want "Error" - %q`, cause, resp, want)
	}
}

func TestInitZeroTimer1(t *testing.T) {
	gid, users := prepareConGame(t, "TestInit(ZeroTimer)", 3)

	rules := c.TestSettings
	rules.Timer1 = 0
	rs, _ := json.Marshal(rules)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  users[0],
		GID:  gid,
	}

	cause, resp := c.Handle(gMsg)
	if cause != "" {
		t.Fatalf(`TestInit(ZeroTimer) = %q - %q, want "" - ""`, cause, resp)
	}
}

func TestNegativeTimer(t *testing.T) {
	gid, users := prepareConGame(t, "TestInit(ZeroTimer)", 3)

	rules := c.TestSettings
	rules.Timer1 = -1
	rs, _ := json.Marshal(rules)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  users[0],
		GID:  gid,
	}

	want := "[Err Msg]"

	cause, resp := c.Handle(gMsg)
	if cause != "Error" {
		t.Fatalf(`TestInit(ZeroTimer) = %q - %q, want "Error" - %q`, cause, resp, want)
	}
}

func TestInitZeroTimer2(t *testing.T) {
	gid, users := prepareConGame(t, "TestInit(ZeroTimer)", 3)

	rules := c.TestSettings
	rules.Timer2 = 0
	rs, _ := json.Marshal(rules)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  users[0],
		GID:  gid,
	}

	cause, resp := c.Handle(gMsg)
	if cause != "" {
		t.Fatalf(`TestInit(ZeroTimer) = %q - %q, want "" - ""`, cause, resp)
	}
}

func TestInitEmptyPrompts(t *testing.T) {
	gid, users := prepareConGame(t, "TestInit(EmptyPrompts)", 3)

	rules := c.TestSettings
	rules.Prompts = []string{}
	rs, _ := json.Marshal(rules)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  users[0],
		GID:  gid,
	}

	cause, resp := c.Handle(gMsg)
	if cause != "" {
		t.Fatalf(`TestInit(EmptyPrompts) = %q - %q want "" - ""`, cause, resp)
	}
}

func TestInitSinglePrompt(t *testing.T) {
	gid, users := prepareConGame(t, "TestInit(SinglePrompt)", 3)

	rules := c.TestSettings
	rules.Prompts = []string{"Single Prompt"}
	rs, _ := json.Marshal(rules)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  users[0],
		GID:  gid,
	}

	want := "[Err Msg]"

	cause, resp := c.Handle(gMsg)
	if cause != "Error" {
		t.Fatalf(`TestInit(SinglePrompts) = %q - %q, want "Error" - %q`, cause, resp, want)
	}
}

func TestInitInvalidStruct(t *testing.T) {
	gid, users := prepareConGame(t, "TestInit(InvalidStruct)", 3)

	rs, _ := json.Marshal(c.DefPlr)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = string(rs)
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  users[0],
		GID:  gid,
	}
	want := "[Err Msg]"

	cause, resp := c.Handle(gMsg)
	if cause != "Error" {
		t.Fatalf(`TestInit(InvalidStruct) = %q - %q, want "Error" - %q`, cause, resp, want)
	}
}

func TestInitEmpty(t *testing.T) {
	gid, users := prepareConGame(t, "TestInit(Empty)", 3)

	game, _ := g.GetGame(gid)
	game.AdditionalRules = ""
	g.UpdateGame(game)

	gMsg := utils.GameMsg{
		Type: "Init",
		UID:  users[0],
		GID:  gid,
	}

	cause, resp := c.Handle(gMsg)
	if cause != "" {
		t.Fatalf(`TestInit(Empty) = %q - %q, want "" - ""`, cause, resp)
	}
}

func TestStart(t *testing.T) {
	gid, users := initConGame(t, "TestStart()", 3)

	gMsg := utils.GameMsg{
		Type: "Start",
		UID:  users[0],
		GID:  gid,
	}

	cause, resp := c.Handle(gMsg)
	if cause != "" {
		t.Fatalf(`TestStart() = %q - %q, want "" - ""`, cause, resp)
	}

	timer := c.TestSettings.Timer1
	want := createWant(c.PROMPTS, users, timer)

	cVars, err := c.GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestPause() = %v, "%v",want %v, "nil"`, cVars, err, want)
	}

}

func TestPauseTimer(t *testing.T) {
	gid, users := startConGame(t, "TestPause()", 3)

	gMsg := utils.GameMsg{
		Type: "Pause",
		UID:  users[0],
		GID:  gid,
	}

	cause, resp := c.Handle(gMsg)
	if cause != "" {
		t.Fatalf(`TestPause() = %q - %q, want "" - ""`, cause, resp)
	}

	want := createWant(c.PROMPTS, users, c.TestSettings.Timer1)
	want.Paused = true

	cVars, err := c.GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestPause() = %v, "%v", want %v, "nil"`, cVars, nil, want)
	}

	checkTimeout(t, "TestPause()", c.TestSettings.Timer1, gid, false)
}

func TestUnpauseTimer(t *testing.T) {
	gid, users := startConGame(t, "TestUnpause()", 3)

	gMsg := utils.GameMsg{
		Type: "Pause",
		UID:  users[0],
		GID:  gid,
	}

	_, _ = c.Handle(gMsg)

	cause, resp := c.Handle(gMsg)
	if cause != "" {
		t.Fatalf(`TestUnpause() = %q - %q, want "" - ""`, cause, resp)
	}

	want := createWant(c.PROMPTS, users, c.TestSettings.Timer1)

	cVars, err := c.GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestUnpause() = %v, "%v", want %v, "nil"`, cVars, nil, want)
	}

	checkTimeout(t, "TestPause()", c.TestSettings.Timer1, gid, true)
}

func TestEnd(t *testing.T) {
	gid, users := startConGame(t, "TestEnd()", 3)

	gMsg := utils.GameMsg{
		Type: "End",
		UID:  users[0],
		GID:  gid,
	}

	cause, resp := c.Handle(gMsg)
	if cause != "" {
		t.Fatalf(`TestEnd() = %q - %q, want "" - ""`, cause, resp)
	}

	want := c.ConVars{}

	cVars, err := c.GetConState(gid)
	if !cmp.Equal(cVars, want) || err == nil {
		t.Fatalf(`TestEnd() = %v, "%v", want %v, ERROR`, cVars, nil, want)
	}
}

func TestReset(t *testing.T) {
	gid, users := startConGame(t, "TestReset()", 3)

	gMsg := utils.GameMsg{
		Type: "Reset",
		UID:  users[0],
		GID:  gid,
	}

	cause, resp := c.Handle(gMsg)
	if cause != "" {
		t.Fatalf(`TestReset() = %q - %q, want "" - ""`, cause, resp)
	}

	want := createWant(c.LOBBY, nil, c.TestSettings.Timer1)

	cVars, err := c.GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestReset() = %v, "%v", want %v, "nil"`, cVars, nil, want)
	}

}

func TestRemove(t *testing.T) {
	gid, users := startConGame(t, "TestRemove()", 3)

	gMsg := utils.GameMsg{
		Type:    "Remove",
		UID:     users[0],
		GID:     gid,
		Content: users[0],
	}

	cause, resp := c.Handle(gMsg)
	if cause != "" {
		t.Fatalf(`TestRemove() = %q - %q, want "" - ""`, cause, resp)
	}

	users[0] = users[len(users)-1]

	want := createWant(c.PROMPTS, users, c.TestSettings.Timer1)

	cVars, err := c.GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestRemove() = %v, "%v", want %v, "nil"`, cVars, nil, want)
	}
}

// Don't need to test remove too few, handled by leave below

func TestStatus(t *testing.T) {
	gid, users := startConGame(t, "TestStatus(ConSpec)", 3)

	status := "Ready"

	gMsg := utils.GameMsg{
		Type:    "Status",
		UID:     users[0],
		GID:     gid,
		Content: status,
	}

	cause, resp := c.Handle(gMsg)
	if cause != "" {
		t.Fatalf(`TestStatus(ConSpec) = %q - %q, want "" - ""`, cause, resp)
	}
}

func TestStatusPhaseChange(t *testing.T) {
	gid, users := startConGame(t, "TestStatus(PhaseChange)", 3)

	gMsg := utils.GameMsg{
		Type:    "Status",
		UID:     users[0],
		GID:     gid,
		Content: "Ready",
	}

	cVar := c.CVars[gid]
	cVar.State = c.STORIES
	c.CVars[gid] = cVar

	for i, user := range users {
		gMsg.UID = user
		cause, resp := c.Handle(gMsg)
		if cause != "" {
			t.Fatalf(`TestStatus(PhaseChange) [%d] = %q - %q, want "" - ""`, i, cause, resp)
		}
	}

	wantVars := createWant(c.PROMPTS, users, c.TestSettings.Timer2)
	wantVars.Round++
	gMsg.UID = users[0]

	cause, resp := c.Handle(gMsg)
	if cause != "" {
		t.Fatalf(`TestStatus(PhaseChange) = %q - %q, want "" - ""`, cause, resp)
	}

	cVars, err := c.GetConState(gid)
	if !cmp.Equal(cVars, wantVars) || err != nil {
		t.Fatalf(`TestStatus(PhaseChange) = %v - %q, want %v, "nil"`, cVars, err, wantVars)
	}
}

func TestLeaveValid(t *testing.T) {
	gid, users := startConGame(t, "TestLeave(Valid)", 3)

	gMsg := utils.GameMsg{
		Type: "Leave",
		UID:  users[0],
		GID:  gid,
	}

	want := createWant(c.PROMPTS, users, c.TestSettings.Timer1)
	delete(want.Stories, users[0])

	cause, resp := c.Handle(gMsg)
	if cause != "" {
		t.Fatalf(`TestLeave(Valid) = %q - %q, want "" - ""`, cause, resp)
	}

	cVars, err := c.GetConState(gid)
	if len(cVars.Stories) != len(want.Stories) || err != nil {
		t.Fatalf(`TestLeave(Valid) = %q, "%v", want %q, "nil"`, cVars.Stories, err, want.Stories)
	}
}

func TestReplyValid(t *testing.T) {
	gid, users := startConGame(t, "TestReply(Valid)", 2)

	story, _ := json.Marshal(c.DefStory)

	gMsg := utils.GameMsg{
		Type:    "Reply",
		UID:     users[0],
		GID:     gid,
		Content: string(story),
	}

	cause, resp := c.Handle(gMsg)
	if cause != "" {
		t.Fatalf(`TestReply(Valid) = %q - %q, want "" - ""`, cause, resp)
	}

	want := createWant(c.PROMPTS, users, c.TestSettings.Timer1)
	want.Stories[users[0]] = c.DefStory

	cVars, err := c.GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestReply(Valid) = %v, "%v", want %v, "nil"`, cVars, err, want)
	}
}

func TestReplyPhaseChange(t *testing.T) {
	gid, users := startConGame(t, "TestReply(PhaseChange)", 2)

	story, _ := json.Marshal(c.DefStory)

	gMsg := utils.GameMsg{
		Type:    "Reply",
		GID:     gid,
		Content: string(story),
	}

	var cause, resp string

	want := createWant(c.POSTPROMPTS, users, c.TestSettings.Timer1)
	for i, user := range users {
		gMsg.UID = user
		cause, resp := c.Handle(gMsg)
		if cause != "" {
			t.Fatalf(`TestReply(PhaseChange) [%d] = %q - %q, want "" - ""`, i, cause, resp)
		}
		want.Stories[user] = c.DefStory
	}

	if cause != "" {
		t.Fatalf(`TestReply(PhaseChange) = %q - %q, want "" - ""`, cause, resp)
	}

	cVars, err := c.GetConState(gid)
	if !cmp.Equal(cVars, want) || err != nil {
		t.Fatalf(`TestReply(Valid) = %v, "%v", want %v, "nil"`, cVars, err, want)
	}

}

func TestReplyShort(t *testing.T) {
	gid, users := startConGame(t, "TestReply(Short)", 2)

	wrongStory := c.DefStory[0:3]
	story, _ := json.Marshal(wrongStory)

	gMsg := utils.GameMsg{
		Type:    "Reply",
		UID:     users[0],
		GID:     gid,
		Content: string(story),
	}

	want := "[Err Msg]"

	cause, resp := c.Handle(gMsg)
	if cause != "Error" {
		t.Fatalf(`TestReply(Short) = %q - %q, want "Error" - %q`, cause, resp, want)
	}

	wantVars := createWant(c.PROMPTS, users, c.TestSettings.Timer1)

	cVars, err := c.GetConState(gid)
	if !cmp.Equal(cVars, wantVars) || err != nil {
		t.Fatalf(`TestReply(Short) = %v, "%v", want %v, "nil"`, cVars, err, wantVars)
	}
}

func TestReplyLong(t *testing.T) {
	gid, users := startConGame(t, "TestReply(Long)", 2)

	wrongStory := append(c.DefStory, "Invalid")
	story, _ := json.Marshal(wrongStory)

	gMsg := utils.GameMsg{
		Type:    "Reply",
		UID:     users[0],
		GID:     gid,
		Content: string(story),
	}

	want := "[Err Msg]"

	cause, resp := c.Handle(gMsg)
	if cause != "Error" {
		t.Fatalf(`TestReply(Long) = %q - %q, want "Error" - %q`, cause, resp, want)
	}

	wantVars := createWant(c.PROMPTS, users, c.TestSettings.Timer1)

	cVars, err := c.GetConState(gid)
	if !cmp.Equal(cVars, wantVars) || err != nil {
		t.Fatalf(`TestReply(Long) = %v, "%v", want %v, "nil"`, cVars, err, wantVars)
	}
}

func TestReplyEmptyLine(t *testing.T) {
	gid, users := startConGame(t, "TestReply(EmptyLine)", 2)

	wrongStory := c.DefStory
	wrongStory[2] = ""
	story, _ := json.Marshal(wrongStory)

	gMsg := utils.GameMsg{
		Type:    "Reply",
		UID:     users[0],
		GID:     gid,
		Content: string(story),
	}

	want := "[Err Msg]"

	cause, resp := c.Handle(gMsg)
	if cause != "Error" {
		t.Fatalf(`TestReply(EmptyLine) = %q - %q, want "Error" - %q`, cause, resp, want)
	}

	wantVars := createWant(c.PROMPTS, users, c.TestSettings.Timer1)

	cVars, err := c.GetConState(gid)
	if !cmp.Equal(cVars, wantVars) || err != nil {
		t.Fatalf(`TestReply(EmptyLine) = %v, "%v", want %v, "nil"`, cVars, err, wantVars)
	}
}

func TestReplyDuplicate(t *testing.T) {
	gid, users := startConGame(t, "TestReply(Duplicate)", 2)

	story, _ := json.Marshal(c.DefStory)

	gMsg := utils.GameMsg{
		Type:    "Reply",
		UID:     users[0],
		GID:     gid,
		Content: string(story),
	}

	_, _ = c.Handle(gMsg)

	want := "[Err Msg]"

	cause, resp := c.Handle(gMsg)
	if cause != "Error" {
		t.Fatalf(`TestReply(Valid) = %q - %q, want "Error" - %q`, cause, resp, want)
	}

	wantVars := createWant(c.PROMPTS, users, c.TestSettings.Timer1)
	wantVars.Stories[users[0]] = c.DefStory

	cVars, err := c.GetConState(gid)
	if !cmp.Equal(cVars, wantVars) || err != nil {
		t.Fatalf(`TestReply(Valid) = %v, "%v", want %v, "nil"`, cVars, err, wantVars)
	}
}
