package user

import (
	"testing"
)

var testUser = user{
	UID:    "",
	Name:   testUserName,
	Status: "New",
}

const testUserName = "Test User"

func TestCreateUser(t *testing.T) {
	id, err := CreateUser(testUserName)
	if id == "" || err != nil {
		t.Fatalf(`CreateUser(%s) = %q, %v, want "uuid", nil`, testUserName, id, err)
	}
}

func TestGetUser(t *testing.T) {
	id, _ := CreateUser(testUserName)

	testUser.UID = id

	user, err := GetUser(id)
	if user != testUser || err != nil {
		t.Fatalf(`GetUser(%s) = %v, %v, want obj, nil`, id, user, err)
	}
}
