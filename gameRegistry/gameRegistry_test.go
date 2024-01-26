package gameRegistry

import (
	"testing"
)

const testAddress = "http://localhost:8091"

const testGameMode = "Test"
const altGameMode = "Alt"
const badGameMode = "Invalid"

func TestRegisterGame(t *testing.T) {
	err := RegisterGameMode(testGameMode, testAddress)
	if err != nil {
		t.Fatalf(`TestRegisterGame(Valid) = %v, want nil`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestRegisterUniqueNameGame(t *testing.T) {
	RegisterGameMode(testGameMode, testAddress)

	err := RegisterGameMode(altGameMode, testAddress)
	if err != nil {
		t.Fatalf(`TestRegisterGame(Unique Name) = %v, want nil`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestRegisterSameNameGame(t *testing.T) {
	RegisterGameMode(testGameMode, testAddress)
	err := RegisterGameMode(testGameMode, testAddress)
	if err == nil {
		t.Fatalf(`TestRegisterGame(Same Name) = %v, want err`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestRegisterGameEmptyMode(t *testing.T) {
	err := RegisterGameMode("", testAddress)
	if err == nil {
		t.Fatalf(`TestRegisterGame(Empty Name) = %v, want err`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestRegisterGameEmptyURL(t *testing.T) {
	err := RegisterGameMode(testGameMode, "")
	if err == nil {
		t.Fatalf(`TestRegisterGame(Empty URL) = %v, want err`, err)
	}
}

func TestRemoveGameFromRegistry(t *testing.T) {
	setupRegisterTest(t)

	err := RemoveGameMode(testGameMode)
	if err != nil {
		t.Fatalf(`TestRemoveGameFromRegistry(Valid) = %v, want nil`, err)
	}
}

func TestRemoveGameFromRegistryInvalidMode(t *testing.T) {
	setupRegisterTest(t)

	err := RemoveGameMode(badGameMode)
	if err == nil {
		t.Fatalf(`TestRemoveGameFromRegistry(Invalid Mode) = %v, want err`, err)
	}
}

func TestRemoveGameFromRegistryDouble(t *testing.T) {
	setupRegisterTest(t)

	RemoveGameMode(testGameMode)
	err := RemoveGameMode(testGameMode)
	if err == nil {
		t.Fatalf(`TestRemoveGameFromRegistry(Double) = %v, want err`, err)
	}
}

func setupRegisterTest(t *testing.T) {
	RegisterGameMode(testGameMode, testAddress)
	RegisterGameMode(altGameMode, testAddress)

	t.Cleanup(cleanUpAfterTest)
}

func cleanUpAfterTest() {
	urlRegistry = make(map[string]string)
}
