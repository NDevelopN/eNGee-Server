package gameRegistry

import "testing"

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

func TestBuildGame(t *testing.T) {
	setupRegisterTest(t)
	addr, err := BuildGame(testGameType)
	if addr != testAddress || err != nil {
		t.Fatalf(`TestBuildGame(Valid) = %q, %v, want %q, nil`, addr, err, testAddress)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestBuildGameMulti(t *testing.T) {
	setupRegisterTest(t)
	BuildGame(testGameType)
	addr, err := BuildGame(testGameType)
	if addr != testAddress || err != nil {
		t.Fatalf(`TestBuildGame(Valid) = %q, %v, want %q, nil`, addr, err, testAddress)
	}
}

func TestBuildGameEmptyType(t *testing.T) {
	setupRegisterTest(t)
	addr, err := BuildGame("")
	if addr != "" || err == nil {
		t.Fatalf(`TestBuildGame(EmptyType) = %q, %v, want "", err`, addr, err)
	}
}

func TestBuildGameInvalidType(t *testing.T) {
	setupRegisterTest(t)
	addr, err := BuildGame(badGameType)
	if addr != "" || err == nil {
		t.Fatalf(`TestBuildGame(InvalidType) = %q, %v, want "", err`, addr, err)
	}
}

func setupRegisterTest(t *testing.T) {
	RegisterGameType(testGameType, testDummyFunc)

	t.Cleanup(cleanUpAfterTest)
}

func cleanUpAfterTest() {
	registry = make(map[string]func() (string, error))
}
