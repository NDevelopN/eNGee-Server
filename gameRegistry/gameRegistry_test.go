package gameRegistry

import (
	"testing"
)

const testAddress = "Address"

const testGameType = "Test"
const altGameType = "Alt"
const badGameType = "Invalid"

func TestRegisterGame(t *testing.T) {
	err := RegisterGameType(testGameType, testAddress)
	if err != nil {
		t.Fatalf(`TestRegisterGame(Valid) = %v, want nil`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestRegisterUniqueNameGame(t *testing.T) {
	RegisterGameType(testGameType, testAddress)

	err := RegisterGameType(altGameType, testAddress)
	if err != nil {
		t.Fatalf(`TestRegisterGame(Unique Name) = %v, want nil`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestRegisterSameNameGame(t *testing.T) {
	RegisterGameType(testGameType, testAddress)
	err := RegisterGameType(testGameType, testAddress)
	if err == nil {
		t.Fatalf(`TestRegisterGame(Same Name) = %v, want err`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestRegisterGameEmptyType(t *testing.T) {
	err := RegisterGameType("", testAddress)
	if err == nil {
		t.Fatalf(`TestRegisterGame(Empty Name) = %v, want err`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestRegisterGameEmptyURL(t *testing.T) {
	err := RegisterGameType(testGameType, "")
	if err == nil {
		t.Fatalf(`TestRegisterGame(Empty URL) = %v, want err`, err)
	}
}

func TestRemoveGameFromRegistry(t *testing.T) {
	setupRegisterTest(t)

	err := RemoveGameType(testGameType)
	if err != nil {
		t.Fatalf(`TestRemoveGameFromRegistry(Valid) = %v, want nil`, err)
	}
}

func TestRemoveGameFromRegistryInvalidType(t *testing.T) {
	setupRegisterTest(t)

	err := RemoveGameType(badGameType)
	if err == nil {
		t.Fatalf(`TestRemoveGameFromRegistry(Invalid Type) = %v, want err`, err)
	}
}

func TestRemoveGameFromRegistryDouble(t *testing.T) {
	setupRegisterTest(t)

	RemoveGameType(testGameType)
	err := RemoveGameType(testGameType)
	if err == nil {
		t.Fatalf(`TestRemoveGameFromRegistry(Double) = %v, want err`, err)
	}
}

func setupRegisterTest(t *testing.T) {
	RegisterGameType(testGameType, testAddress)
	RegisterGameType(altGameType, testAddress)

	t.Cleanup(cleanUpAfterTest)
}

func cleanUpAfterTest() {
	urlRegistry = make(map[string]string)
}
