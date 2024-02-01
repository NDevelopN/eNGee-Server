package gameRegistry

import (
	"errors"
	"testing"

	sErr "Engee-Server/stockErrors"
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
	if !errors.As(err, &sErr.MF_ERR) {
		t.Fatalf(`TestRegisterGame(Same Name) = %v, want MatchFoundError`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestRegisterGameEmptyMode(t *testing.T) {
	err := RegisterGameMode("", testAddress)
	if !errors.As(err, &sErr.EV_ERR) {
		t.Fatalf(`TestRegisterGame(Empty Name) = %v, want EmptyValueError`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestRegisterGameEmptyURL(t *testing.T) {
	err := RegisterGameMode(testGameMode, "")
	if !errors.As(err, &sErr.EV_ERR) {
		t.Fatalf(`TestRegisterGame(Empty URL) = %v, want EmptyValueError`, err)
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
	if !errors.As(err, &sErr.MNF_ERR) {
		t.Fatalf(`TestRemoveGameFromRegistry(Invalid Mode) = %v, want MatchNotFoundError`, err)
	}
}

func TestRemoveGameFromRegistryDouble(t *testing.T) {
	setupRegisterTest(t)

	RemoveGameMode(testGameMode)
	err := RemoveGameMode(testGameMode)
	if !errors.As(err, &sErr.MNF_ERR) {
		t.Fatalf(`TestRemoveGameFromRegistry(Double) = %v, want MatchNotFoundError`, err)
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
