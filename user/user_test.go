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
const newTestName = "New Name"

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

func TestUpdateUserName(t *testing.T) {
	id, _ := CreateUser(testUserName)

	tuInstance := testUser
	tuInstance.UID = id
	tuInstance.Name = newTestName

	err := UpdateUserName(id, newTestName)
	if err != nil {
		t.Fatalf(`UpdateUserName(%s, %s) = %v, want nil`, id, newTestName, err)
	}

	user, err := GetUser(id)
	if user != tuInstance || err != nil {
		t.Fatalf(`UpdateUserName(%s, %s) = %v, %v, want %v, nil`, id, newTestName, user, err, tuInstance)
	}
}
