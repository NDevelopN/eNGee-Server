package gameRegistry

import (
	"testing"
)

const testAddress = "Address"

const testGameType = "Test"
const altGameType = "Alt"
const badGameType = "Invalid"

var testDummyFunc = func() (string, error) {
	return testAddress, nil
}

func TestRegisterGame(t *testing.T) {
	err := RegisterGameType(testGameType, testDummyFunc)
	if err != nil {
		t.Fatalf(`TestRegisterGame(Valid) = %v, want nil`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestRegisterUniqueNameGame(t *testing.T) {
	RegisterGameType(testGameType, testDummyFunc)

	err := RegisterGameType(altGameType, testDummyFunc)
	if err != nil {
		t.Fatalf(`TestRegisterGame(Unique Name) = %v, want nil`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestRegisterSameNameGame(t *testing.T) {
	RegisterGameType(testGameType, testDummyFunc)
	err := RegisterGameType(testGameType, testDummyFunc)
	if err == nil {
		t.Fatalf(`TestRegisterGame(Same Name) = %v, want err`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestRegisterGameEmptyType(t *testing.T) {
	err := RegisterGameType("", testDummyFunc)
	if err == nil {
		t.Fatalf(`TestRegisterGame(Empty Name) = %v, want err`, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestRemoveGameFromRegistry(t *testing.T) {
	setupRegisterTest(t)

	err := RemoveGame(testGameType)
	if err != nil {
		t.Fatalf(`TestRemoveGameFromRegistry(Valid) = %v, want nil`, err)
	}
}

func TestRemoveGameFromRegistryInvalidType(t *testing.T) {
	setupRegisterTest(t)

	err := RemoveGame(badGameType)
	if err == nil {
		t.Fatalf(`TestRemoveGameFromRegistry(Invalid Type) = %v, want err`, err)
	}
}

func TestRemoveGameFromRegistryDouble(t *testing.T) {
	setupRegisterTest(t)

	RemoveGame(testGameType)
	err := RemoveGame(testGameType)
	if err == nil {
		t.Fatalf(`TestRemoveGameFromRegistry(Double) = %v, want err`, err)
	}
}

func setupRegisterTest(t *testing.T) {
	RegisterGameType(testGameType, testDummyFunc)

	t.Cleanup(cleanUpAfterTest)
}

func cleanUpAfterTest() {
	registry = make(map[string]func() (string, error))
}
