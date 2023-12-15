package gamedummy

import (
	"testing"
)

var defaultGame = GameDummy{
	Rules:  dummyRules,
	Status: NEW,
}

var emptyGame = GameDummy{}

func TestCreateDefaultGame(t *testing.T) {
	createdGame := CreateDefaultGame()
	if createdGame != defaultGame {
		t.Fatalf(`TestCreateDefaultGame() = %v, want %v`, createdGame, defaultGame)
	}
}

func TestStartGame(t *testing.T) {
	testGame := defaultGame

	expected := defaultGame
	expected.Status = ACTIVE

	err := testGame.StartGame()
	if testGame != expected || err != nil {
		t.Fatalf(`TestStartGame(Valid) = %v, %v, want %v, nil`, testGame, err, expected)
	}
}

func TestStartGameEmptyGame(t *testing.T) {
	testGame := emptyGame
	expected := emptyGame

	err := testGame.StartGame()
	if testGame != expected || err == nil {
		t.Fatalf(`TestStartGame(EmptyGame) = %v, %v, want %v, err`, testGame, err, expected)
	}
}

func TestStartGameInvalidStatus(t *testing.T) {
	testGame := defaultGame
	expected := defaultGame

	testGame.Status = ACTIVE
	expected.Status = ACTIVE

	err := testGame.StartGame()
	if testGame != expected || err == nil {
		t.Fatalf(`TestStartGame(InvalidStatus) = %v, %v, want %v, err`, testGame, err, expected)
	}
}

func TestStartGameRESET(t *testing.T) {
	testGame := defaultGame
	expected := defaultGame

	testGame.Status = RESET
	expected.Status = ACTIVE

	err := testGame.StartGame()
	if testGame != expected || err != nil {
		t.Fatalf(`TestStartGame(RESET) = %v, %v, want %v, nil`, testGame, err, expected)
	}
}

func TestPauseGame(t *testing.T) {
	testGame := setUpTestGame()

	expected := testGame
	expected.Status = PAUSED

	err := testGame.PauseGame()
	if testGame != expected || err != nil {
		t.Fatalf(`TestPauseGame(Valid) = %v, %v, want %v, nil`, testGame, err, expected)
	}
}

func TestPauseGameEmptyGame(t *testing.T) {
	testGame := emptyGame
	expected := emptyGame

	err := testGame.PauseGame()
	if testGame != expected || err == nil {
		t.Fatalf(`TestPauseGame(EmptyGame) = %v, %v, want %v, err`, testGame, err, expected)
	}
}

func TestPauseGameInvalidStatus(t *testing.T) {
	testGame := defaultGame
	expected := defaultGame

	err := testGame.PauseGame()
	if testGame != expected || err == nil {
		t.Fatalf(`TestPauseGame(InvalidStatus) = %v, %v, want %v, err`, testGame, err, expected)
	}
}

func TestPauseGamePAUSED(t *testing.T) {
	testGame := setUpTestGame()
	expected := testGame

	testGame.Status = PAUSED
	expected.Status = ACTIVE

	err := testGame.PauseGame()
	if testGame != expected || err != nil {
		t.Fatalf(`TestPauseGame(PAUSED) = %v, %v, want %v, nil`, testGame, err, expected)
	}
}

func TestResetGame(t *testing.T) {
	testGame := setUpTestGame()

	expected := testGame
	expected.Status = RESET

	err := testGame.ResetGame()
	if testGame != expected || err != nil {
		t.Fatalf(`TestResetGame(Valid) = %v, %v, want %v, nil`, testGame, err, expected)
	}
}

func TestResetGameEmptyGame(t *testing.T) {
	testGame := emptyGame
	expected := emptyGame

	err := testGame.ResetGame()
	if testGame != expected || err == nil {
		t.Fatalf(`TestResetGame(EmptyGame) = %v, %v, want %v, err`, testGame, err, emptyGame)
	}
}

func TestResetGameInvalidStatus(t *testing.T) {
	testGame := defaultGame
	expected := defaultGame

	err := testGame.ResetGame()
	if testGame != expected || err == nil {
		t.Fatalf(`TestResetGame(InvalidStatus) = %v, %v, want %v, err`, testGame, err, expected)
	}
}

func TestResetGamePAUSED(t *testing.T) {
	testGame := setUpTestGame()
	expected := testGame

	testGame.Status = PAUSED
	expected.Status = RESET

	err := testGame.ResetGame()
	if testGame != expected || err != nil {
		t.Fatalf(`TestResetGame(PAUSED) = %v, %v, want %v, nil`, testGame, err, expected)
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
	testGame := emptyGame

	err := testGame.EndGame()
	if err == nil {
		t.Fatalf(`TestEndGame(EmptyGame) = %v, want err`, err)
	}
}

func TestEndGamePAUSED(t *testing.T) {
	testGame := setUpTestGame()
	testGame.Status = PAUSED

	err := testGame.EndGame()
	if err != nil {
		t.Fatalf(`TestEndGame(PAUSED) = %v, want nil`, err)
	}
}

func setUpTestGame() GameDummy {
	testGame := defaultGame
	testGame.StartGame()

	return testGame
}
