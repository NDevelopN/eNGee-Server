package gameRegistry

import "testing"

const testAddress = "Address"
const testGameType = "Test"
const altGameType = "Alt"

var testDummyFunc = func() (string, error) {
	return testAddress, nil
}

func TestRegisterGame(t *testing.T) {
	err := RegisterGameType(testGameType, testDummyFunc)
	if err != nil {
		t.Fatalf(`TestRegisterGame(Valid) = %v, want nil`, err)
	}
}

func TestRegisterUniqueNameGame(t *testing.T) {
	RegisterGameType(testGameType, testDummyFunc)

	err := RegisterGameType(altGameType, testDummyFunc)
	if err != nil {
		t.Fatalf(`TestRegisterGame(Unique Name) = %v, want nil`, err)
	}
}

func TestRegisterSameNameGame(t *testing.T) {
	RegisterGameType(testGameType, testDummyFunc)
	err := RegisterGameType(testGameType, testDummyFunc)
	if err == nil {
		t.Fatalf(`TestRegisterGame(Same Name) = %v, want err`, err)
	}
}

func TestRegisterGameEmptyType(t *testing.T) {
	err := RegisterGameType("", testDummyFunc)
	if err == nil {
		t.Fatalf(`TestRegisterGame(Empty Name) = %v, want err`, err)
	}
}
