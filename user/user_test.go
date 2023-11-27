package user

import (
	"testing"
)

const testUserName = "Test User"

func TestCreateUser(t *testing.T) {
	id, err := CreateUser(testUserName)
	if id == "" || err != nil {
		t.Fatalf(`CreateUser() = %q, "%v", want "uuid", "nil"`, id, err)
	}
}
