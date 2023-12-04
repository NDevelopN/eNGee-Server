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

const testConPort = "8091"
const altConPort = "8092"
const testConURL = "http://localhost:" + testConPort
const altConURL = "http://localhost:" + altConPort
const testPlayURL = "http://localhost:8099"

const badURL = "http://notahost:8080"

func TestMain(m *testing.M) {
	setupGameSuite()
	code := m.Run()
	cleanUpAfterSuite()
	os.Exit(code)
}

func TestCreateGame(t *testing.T) {
	pURL, err := CreateGame(testRID, testConURL)
	if pURL != testPlayURL || err != nil {
		t.Fatalf(`TestCreateGame(Valid) = %q, %v, want %q, nil`, pURL, err, testPlayURL)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameDoubleSameURL(t *testing.T) {
	CreateGame(testRID, testConURL)
	pURL, err := CreateGame(testRID, testConURL)
	if pURL != "" || err == nil {
		t.Fatalf(`TestCreateGame(Double Same) = %q, %v, want "", err`, pURL, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameDoubleUniqueURL(t *testing.T) {
	CreateGame(testRID, testConURL)
	pURL, err := CreateGame(testRID, altConURL)
	if pURL != "" || err == nil {
		t.Fatalf(`TestCreateGame(Double Unique) = %q, %v, want "", err`, pURL, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameMultiSameURL(t *testing.T) {
	CreateGame(testRID, testConURL)
	pURL, err := CreateGame(altRID, testConURL)
	if pURL != testPlayURL || err != nil {
		t.Fatalf(`TestCreateGame(Same URL) = %q, %v, want %q, nil`, pURL, err, testPlayURL)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameMultiUniqueURL(t *testing.T) {
	CreateGame(testRID, testConURL)
	pURL, err := CreateGame(altRID, altConURL)
	if pURL != testPlayURL || err != nil {
		t.Fatalf(`TestCreateGame(Unique URL) = %q, %v, want %q, nil`, pURL, err, testPlayURL)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameEmptyRID(t *testing.T) {
	pURL, err := CreateGame("", testConURL)
	if pURL != "" || err == nil {
		t.Fatalf(`TestCreateGame(Empty RID) = %q, %v, want "", err`, pURL, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameEmptyURL(t *testing.T) {
	pURL, err := CreateGame(testRID, "")
	if pURL != "" || err == nil {
		t.Fatalf(`TestCreateGame(Empty URL) = %q, %v, want "", err`, pURL, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateGameInvalidURL(t *testing.T) {
	pURL, err := CreateGame(testRID, badURL)
	if pURL != "" || err == nil {
		t.Fatalf(`TestCreateGame(Valid) = %q, %v, want "", err`, pURL, err)
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

func setupGameSuite() {
	go gamedummy.Start(testConPort)
	go gamedummy.Start(altConPort)

	time.Sleep(200 * time.Millisecond)
}

func setupGameTest(t *testing.T) {

	CreateGame(testRID, testConURL)
	CreateGame(altRID, altConURL)

	t.Cleanup(cleanUpAfterTest)
}

func cleanUpAfterTest() {
	EndGame(testRID)
	EndGame(altRID)

	gameURLs = make(map[string]string)
}

func cleanUpAfterSuite() {

}
