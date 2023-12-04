package gamedummy

import (
	"testing"
)

var defaultGame = GameDummy{
	Address: dummyAddress,
	Rules:   dummyRules,
	Status:  NEW,
}

var emptyGame = GameDummy{}

func TestCreateDefaultGame(t *testing.T) {
	createdGame := CreateDefaultGame()
	if createdGame != defaultGame {
		t.Fatalf(`TestCreateDefaultGame() = %v, want %v`, createdGame, defaultGame)
	}
}

func TestStartGame(t *testing.T) {
	expected := defaultGame
	expected.Status = ACTIVE

	startedGame, err := defaultGame.StartGame()
	if startedGame != expected || err != nil {
		t.Fatalf(`TestStartGame(Valid) = %v, %v, want %v, nil`, startedGame, err, expected)
	}
}

func TestStartGameEmptyGame(t *testing.T) {
	startedGame, err := emptyGame.StartGame()
	if startedGame != emptyGame || err == nil {
		t.Fatalf(`TestStartGame(EmptyGame) = %v, %v, want %v, err`, startedGame, err, emptyGame)
	}
}

func TestStartGameInvalidStatus(t *testing.T) {
	testGame := defaultGame
	testGame.Status = ACTIVE

	startedGame, err := testGame.StartGame()
	if startedGame != testGame || err == nil {
		t.Fatalf(`TestStartGame(InvalidStatus) = %v, %v, want %v, err`, startedGame, err, testGame)
	}
}

func TestStartGameRESET(t *testing.T) {
	testGame := defaultGame
	testGame.Status = RESET

	expected := defaultGame
	expected.Status = ACTIVE

	startedGame, err := testGame.StartGame()
	if startedGame != expected || err != nil {
		t.Fatalf(`TestStartGame(RESET) = %v, %v, want %v, nil`, startedGame, err, expected)
	}
}

func TestPauseGame(t *testing.T) {
	testGame := setUpTestGame()

	expected := testGame
	expected.Status = PAUSED

	pausedGame, err := testGame.PauseGame()
	if pausedGame != expected || err != nil {
		t.Fatalf(`TestPauseGame(Valid) = %v, %v, want %v, nil`, pausedGame, err, expected)
	}
}

func TestPauseGameEmptyGame(t *testing.T) {
	pausedGame, err := emptyGame.PauseGame()
	if pausedGame != emptyGame || err == nil {
		t.Fatalf(`TestPauseGame(EmptyGame) = %v, %v, want %v, err`, pausedGame, err, emptyGame)
	}
}

func TestPauseGameInvalidStatus(t *testing.T) {
	pausedGame, err := defaultGame.PauseGame()
	if pausedGame != defaultGame || err == nil {
		t.Fatalf(`TestPauseGame(InvalidStatus) = %v, %v, want %v, err`, pausedGame, err, defaultGame)
	}
}

func TestPauseGamePAUSED(t *testing.T) {
	testGame := setUpTestGame()
	testGame, _ = testGame.PauseGame()

	expected := testGame
	expected.Status = ACTIVE

	pausedGame, err := testGame.PauseGame()
	if pausedGame != expected || err != nil {
		t.Fatalf(`TestPauseGame(PAUSED) = %v, %v, want %v, nil`, pausedGame, err, expected)
	}
}

func TestResetGame(t *testing.T) {
	testGame := setUpTestGame()

	expected := testGame
	expected.Status = RESET

	resetGame, err := testGame.ResetGame()
	if resetGame != expected || err != nil {
		t.Fatalf(`TestResetGame(Valid) = %v, %v, want %v, nil`, resetGame, err, expected)
	}
}

func TestResetGameEmptyGame(t *testing.T) {
	resetGame, err := emptyGame.ResetGame()
	if resetGame != emptyGame || err == nil {
		t.Fatalf(`TestResetGame(EmptyGame) = %v, %v, want %v, err`, resetGame, err, emptyGame)
	}
}

func TestResetGameInvalidStatus(t *testing.T) {
	resetGame, err := defaultGame.ResetGame()
	if resetGame != defaultGame || err == nil {
		t.Fatalf(`TestResetGame(InvalidStatus) = %v, %v, want %v, err`, resetGame, err, defaultGame)
	}
}

func TestResetGamePAUSED(t *testing.T) {
	testGame := setUpTestGame()
	testGame, _ = testGame.PauseGame()

	expected := testGame
	expected.Status = RESET

	resetGame, err := testGame.ResetGame()
	if resetGame != expected || err != nil {
		t.Fatalf(`TestResetGame(PAUSED) = %v, %v, want %v, nil`, resetGame, err, expected)
	}
}

func TestEndGame(t *testing.T) {
	testGame := setUpTestGame()

	err := testGame.EndGame()
	if err != nil {
		t.Fatalf(`TestEndGame(Valid) = %v, want nil`, err)
	}
}

func TestEndGameEmptyGame(t *testing.T) {
	err := emptyGame.EndGame()
	if err == nil {
		t.Fatalf(`TestEndGame(EmptyGame) = %v, want err`, err)
	}
}

func TestEndGamePAUSED(t *testing.T) {
	testGame := setUpTestGame()
	testGame, _ = testGame.PauseGame()

	err := testGame.EndGame()
	if err != nil {
		t.Fatalf(`TestEndGame(PAUSED) = %v, want nil`, err)
	}
}

func setUpTestGame() GameDummy {
	testGame, _ := defaultGame.StartGame()
	return testGame
}
