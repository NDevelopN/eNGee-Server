package gameclient

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"

	reg "Engee-Server/gameRegistry"
	sErr "Engee-Server/stockErrors"
	"Engee-Server/testDummy"
)

var testRID = uuid.NewString()
var altRID = uuid.NewString()
var badRID = uuid.NewString()

const testGameMode = "Test"
const altGameMode = "Alt"

const testPort = "8091"
const altPort = "8092"
const testURL = "http://localhost:" + testPort
const altURL = "http://localhost:" + altPort

const updatedRules = "New Rules"

const badURL = "http://notahost:8080"

func TestMain(m *testing.M) {
	setupGameSuite()
	code := m.Run()
	cleanUpAfterSuite()
	os.Exit(code)
}

func TestCreateGame(t *testing.T) {
	err := CreateGameInstance(testRID, testURL)
	if err != nil {
		t.Fatalf(`TestCreateGame(Valid) = %v, want nil`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameDoubleSameURL(t *testing.T) {
	CreateGameInstance(testRID, testURL)
	err := CreateGameInstance(testRID, testURL)
	if !errors.As(err, &sErr.MF_ERR) {
		t.Fatalf(`TestCreateGame(Double Same) = %v, want MatchFoundError`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameDoubleUniqueURL(t *testing.T) {
	CreateGameInstance(testRID, testURL)
	err := CreateGameInstance(testRID, altURL)
	if !errors.As(err, &sErr.MF_ERR) {
		t.Fatalf(`TestCreateGame(Double Unique) = %v, want MatchFoundError`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameMultiSameURL(t *testing.T) {
	CreateGameInstance(testRID, testURL)
	err := CreateGameInstance(altRID, testURL)
	if err != nil {
		t.Fatalf(`TestCreateGame(Same URL) = %v, want nil`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameMultiUniqueURL(t *testing.T) {
	CreateGameInstance(testRID, testURL)
	err := CreateGameInstance(altRID, altURL)
	if err != nil {
		t.Fatalf(`TestCreateGame(Unique URL) = %v, want nil`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameEmptyRID(t *testing.T) {
	err := CreateGameInstance("", testURL)
	if !errors.As(err, &sErr.EV_ERR) {
		t.Fatalf(`TestCreateGame(Empty RID) %v, want EmptyValueError`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameEmptyURL(t *testing.T) {
	err := CreateGameInstance(testRID, "")
	if !errors.As(err, &sErr.EV_ERR) {
		t.Fatalf(`TestCreateGame(Empty URL) %v, want EmptyValueError`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameInvalidURL(t *testing.T) {
	err := CreateGameInstance(testRID, badURL)
	if err == nil {
		t.Fatalf(`TestCreateGame(Valid) %v, want error`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestEndGame(t *testing.T) {
	setupGameTest(t)
	err := EndGame(testRID)
	if err != nil {
		t.Fatalf(`TestEndGame(Valid) = %v, want nil`, err)
	}
}
func TestEndGameDouble(t *testing.T) {
	setupGameTest(t)
	EndGame(testRID)
	err := EndGame(testRID)
	if !errors.As(err, &sErr.MNF_ERR) {
		t.Fatalf(`TestEndGame(Double) = %v, want MatchNotFoundError`, err)
	}
}
func TestEndGameMulti(t *testing.T) {
	setupGameTest(t)
	EndGame(testRID)
	err := EndGame(altRID)
	if err != nil {
		t.Fatalf(`TestEndGame(Multi) = %v, want nil`, err)
	}
}
func TestEndGameInvalidRID(t *testing.T) {
	setupGameTest(t)
	err := EndGame(badRID)
	if !errors.As(err, &sErr.MNF_ERR) {
		t.Fatalf(`TestEndGame(InvalidRID) = %v, want MatchNotFoundError`, err)
	}
}
func TestEndGameEmptyRID(t *testing.T) {
	setupGameTest(t)
	err := EndGame("")
	if !errors.As(err, &sErr.EV_ERR) {
		t.Fatalf(`TestEndGame(EmptyRID) = %v, want EmptyValueError`, err)
	}
}

func TestSetGameRules(t *testing.T) {
	setupGameTest(t)

	err := SetGameRules(testRID, updatedRules)
	if err != nil {
		t.Fatalf(`TestSetGameRules(Valid) = %v, want nil`, err)
	}
}

func TestSetGameRulesDouble(t *testing.T) {
	setupGameTest(t)

	SetGameRules(testRID, updatedRules)
	err := SetGameRules(testRID, updatedRules)
	if err != nil {
		t.Fatalf(`TestSetGameRules(Double) = %v, want nil`, err)
	}
}

func TestSetGameRulesInvalidRID(t *testing.T) {
	setupGameTest(t)

	err := SetGameRules(badRID, updatedRules)
	if !errors.As(err, &sErr.MNF_ERR) {
		t.Fatalf(`TestSetGameRules(InvalidRID) = %v, want MatchNotFoundError`, err)
	}
}

func TestStartGame(t *testing.T) {
	setupGameTest(t)

	err := StartGame(testRID)
	if err != nil {
		t.Fatalf(`TestStartGame(Valid) = %v, want nil`, err)
	}
}

func TestStartGameInvalidRID(t *testing.T) {
	setupGameTest(t)

	err := StartGame(badRID)
	if !errors.As(err, &sErr.MNF_ERR) {
		t.Fatalf(`TestStartGame(Invalid RID) = %v, want MatchNotFoundError`, err)
	}
}

func TestPauseGame(t *testing.T) {
	setupActiveGameTest(t)

	err := PauseGame(testRID)
	if err != nil {
		t.Fatalf(`TestPauseGame(Valid) = %v, want nil`, err)
	}
}

func TestPauseGameInvalidRID(t *testing.T) {
	setupActiveGameTest(t)

	err := PauseGame(badRID)
	if !errors.As(err, &sErr.MNF_ERR) {
		t.Fatalf(`TestPauseGame(InvalidRID) = %v, want MatchNotFoundError`, err)
	}
}
func TestResetGame(t *testing.T) {
	setupActiveGameTest(t)

	err := ResetGame(testRID)
	if err != nil {
		t.Fatalf(`TestResetGame(Valid) = %v, want nil`, err)
	}
}

func TestResetGameInvalidRID(t *testing.T) {
	setupActiveGameTest(t)

	err := ResetGame(badRID)
	if !errors.As(err, &sErr.MNF_ERR) {
		t.Fatalf(`TestResetGame(InvalidRID) = %v, want MatchNotFoundError`, err)
	}
}

func setupGameSuite() {
	go testDummy.Serve(testPort)
	go testDummy.Serve(altPort)

	reg.RegisterGameMode(testGameMode, testURL)
	reg.RegisterGameMode(altGameMode, altURL)

	time.Sleep(200 * time.Millisecond)

}

func setupGameTest(t *testing.T) {

	CreateGameInstance(testRID, testURL)
	CreateGameInstance(altRID, altURL)

	t.Cleanup(cleanUpAfterTest)
}

func setupActiveGameTest(t *testing.T) {
	setupGameTest(t)

	StartGame(testRID)
	StartGame(altRID)
}

func cleanUpAfterTest() {
	EndGame(testRID)
	EndGame(altRID)

	gameURLs = make(map[string]string)
}

func cleanUpAfterSuite() {

}
