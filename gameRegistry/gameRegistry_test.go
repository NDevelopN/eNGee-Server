package gameRegistry

import (
	"testing"

	"github.com/google/uuid"
)

const testAddress = "Address"

const testGameType = "Test"
const altGameType = "Alt"
const badGameType = "Invalid"

var testRID = uuid.NewString()
var altRID = uuid.NewString()

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

func TestSelectRoomGame(t *testing.T) {
	setupRegisterTest(t)

	err := SelectRoomGame(testRID, testGameType)
	if err != nil {
		t.Fatalf(`TestSelectRoomGame(Valid) = %v, want nil`, err)
	}
}

func TestSelectRoomChangeSameType(t *testing.T) {
	setupRegisterTest(t)

	SelectRoomGame(testRID, testGameType)
	err := SelectRoomGame(testRID, testGameType)
	if err != nil {
		t.Fatalf(`TestSelectRoomGame(SameType) = %v, want nil`, err)
	}
}

func TestSelectRoomGameChangeUniqueTypes(t *testing.T) {
	setupRegisterTest(t)

	SelectRoomGame(testRID, testGameType)
	err := SelectRoomGame(testRID, altGameType)
	if err != nil {
		t.Fatalf(`TestSelectRoomGame(UniqueTypes) = %v, want nil`, err)
	}
}

func TestSelectRoomGameInvalidType(t *testing.T) {
	setupRegisterTest(t)

	err := SelectRoomGame(testRID, badGameType)
	if err == nil {
		t.Fatalf(`TestSelectRoomGame(InvalidType) = %v, want err`, err)
	}
}

func TestSelectRoomGameChangeInvalidType(t *testing.T) {
	setupRegisterTest(t)

	SelectRoomGame(testRID, testGameType)
	err := SelectRoomGame(testRID, badGameType)
	if err == nil {
		t.Fatalf(`TestSelectRoomGame(ChangeInvalidType) = %v, want err`, err)
	}
}

func TestSelectRoomSameEmptyRID(t *testing.T) {
	setupRegisterTest(t)

	err := SelectRoomGame("", testGameType)
	if err == nil {
		t.Fatalf(`TestSelectRoomGame(EmptyRID) = %v, want err`, err)
	}
}

func TestSelectRoomMultiSame(t *testing.T) {
	setupRegisterTest(t)

	SelectRoomGame(testRID, testGameType)
	err := SelectRoomGame(altRID, testGameType)
	if err != nil {
		t.Fatalf(`TestSelectRoomGame(MultiSame) = %v, want nil`, err)
	}
}

func TestSelectRoomMultiUnique(t *testing.T) {
	setupRegisterTest(t)

	SelectRoomGame(testRID, testGameType)
	err := SelectRoomGame(altRID, altGameType)
	if err != nil {
		t.Fatalf(`TestSelectRoomGame(MultiUnique) = %v, want nil`, err)
	}
}

func setupRegisterTest(t *testing.T) {
	RegisterGameType(testGameType, testDummyFunc)
	RegisterGameType(altGameType, testDummyFunc)

	t.Cleanup(cleanUpAfterTest)
}

func cleanUpAfterTest() {
	registry = make(map[string]func() (string, error))
}
