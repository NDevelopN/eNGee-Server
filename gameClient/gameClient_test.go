package gameclient

import (
	"os"
	"testing"
	"time"

	"Engee-Server/gamedummy"

	"github.com/google/uuid"
)

var testRID = uuid.NewString()
var altRID = uuid.NewString()
var badRID = uuid.NewString()

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
	err := CreateGame(testRID, testURL)
	if err != nil {
		t.Fatalf(`TestCreateGame(Valid) = %v, want nil`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameDoubleSameURL(t *testing.T) {
	CreateGame(testRID, testURL)
	err := CreateGame(testRID, testURL)
	if err == nil {
		t.Fatalf(`TestCreateGame(Double Same) = %v, want err`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameDoubleUniqueURL(t *testing.T) {
	CreateGame(testRID, testURL)
	err := CreateGame(testRID, altURL)
	if err == nil {
		t.Fatalf(`TestCreateGame(Double Unique) = %v, want err`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameMultiSameURL(t *testing.T) {
	CreateGame(testRID, testURL)
	err := CreateGame(altRID, testURL)
	if err != nil {
		t.Fatalf(`TestCreateGame(Same URL) = %v, want nil`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameMultiUniqueURL(t *testing.T) {
	CreateGame(testRID, testURL)
	err := CreateGame(altRID, altURL)
	if err != nil {
		t.Fatalf(`TestCreateGame(Unique URL) = %v, want nil`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameEmptyRID(t *testing.T) {
	err := CreateGame("", testURL)
	if err == nil {
		t.Fatalf(`TestCreateGame(Empty RID) %v, want err`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameEmptyURL(t *testing.T) {
	err := CreateGame(testRID, "")
	if err == nil {
		t.Fatalf(`TestCreateGame(Empty URL) %v, want err`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameInvalidURL(t *testing.T) {
	err := CreateGame(testRID, badURL)
	if err == nil {
		t.Fatalf(`TestCreateGame(Valid) %v, want err`, err)
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
	if err == nil {
		t.Fatalf(`TestEndGame(Double) = %v, want err`, err)
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
	if err == nil {
		t.Fatalf(`TestEndGame(InvalidRID) = %v, want err`, err)
	}
}
func TestEndGameEmptyRID(t *testing.T) {
	setupGameTest(t)
	err := EndGame("")
	if err == nil {
		t.Fatalf(`TestEndGame(EmptyRID) = %v, want err`, err)
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

func TestSetGameRulesAfterStart(t *testing.T) {
	setupActiveGameTest(t)

	err := SetGameRules(testRID, updatedRules)
	if err == nil {
		t.Fatalf(`TestSetGameRules(After Start) = %v, want err`, err)
	}
}

func TestSetGameRulesInvalidRID(t *testing.T) {
	setupGameTest(t)

	err := SetGameRules(badRID, updatedRules)
	if err == nil {
		t.Fatalf(`TestSetGameRules(InvalidRID) = %v, want err`, err)
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
	if err == nil {
		t.Fatalf(`TestStartGame(Invalid RID) = %v, want err`, err)
	}
}

func TestStartGameDouble(t *testing.T) {
	setupGameTest(t)

	StartGame(testRID)
	err := StartGame(testRID)
	if err == nil {
		t.Fatalf(`TestStartGame(Double) = %v, want err`, err)
	}
}

func TestPauseGame(t *testing.T) {
	setupActiveGameTest(t)

	err := PauseGame(testRID)
	if err != nil {
		t.Fatalf(`TestPauseGame(Valid) = %v, want nil`, err)
	}
}

func TestPauseGameDouble(t *testing.T) {
	setupActiveGameTest(t)

	PauseGame(testRID)
	err := PauseGame(testRID)
	if err != nil {
		t.Fatalf(`TestPauseGame(Double) = %v, want nil`, err)
	}
}

func TestPauseGameNotStarted(t *testing.T) {
	setupGameTest(t)

	err := PauseGame(testRID)
	if err == nil {
		t.Fatalf(`TestPauseGame(Not Started) = %v, want err`, err)
	}
}

func TestPauseGameInvalidRID(t *testing.T) {
	setupActiveGameTest(t)

	err := PauseGame(badRID)
	if err == nil {
		t.Fatalf(`TestPauseGame(InvalidRID) = %v, want err`, err)
	}
}
func TestResetGame(t *testing.T) {
	setupActiveGameTest(t)

	err := ResetGame(testRID)
	if err != nil {
		t.Fatalf(`TestResetGame(Valid) = %v, want nil`, err)
	}
}

func TestResetGameDouble(t *testing.T) {
	setupActiveGameTest(t)

	ResetGame(testRID)
	err := ResetGame(testRID)
	if err == nil {
		t.Fatalf(`TestResetGame(Double) = %v, want err`, err)
	}
}

func TestResetGameNotStarted(t *testing.T) {
	setupGameTest(t)

	err := ResetGame(testRID)
	if err == nil {
		t.Fatalf(`TestResetGame(Not Started) = %v, want err`, err)
	}
}

func TestResetGameInvalidRID(t *testing.T) {
	setupActiveGameTest(t)

	err := ResetGame(badRID)
	if err == nil {
		t.Fatalf(`TestResetGame(InvalidRID) = %v, want err`, err)
	}
}

func setupGameSuite() {
	go gamedummy.Start(testPort)
	go gamedummy.Start(altPort)

	time.Sleep(200 * time.Millisecond)
}

func setupGameTest(t *testing.T) {

	CreateGame(testRID, testURL)
	CreateGame(altRID, altURL)

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
