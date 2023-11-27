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

const updatedStatus = "Updated"

func TestCreateUser(t *testing.T) {
	id, err := CreateUser(testUserName)
	if id == "" || err != nil {
		t.Fatalf(`CreateUser(%s) = %q, %v, want "uuid", nil`, testUserName, id, err)
	}
}

func TestGetUser(t *testing.T) {
	id, tuInstance := setupUserTest()

	user, err := GetUser(id)
	if user != tuInstance || err != nil {
		t.Fatalf(`GetUser(%s) = %v, %v, want obj, nil`, id, user, err)
	}
}

func TestUpdateUserName(t *testing.T) {
	id, tuInstance := setupUserTest()

	tuInstance.Name = newTestName

	err := UpdateUserName(id, newTestName)
	if err != nil {
		t.Fatalf(`UpdateUserName(%s, %s) = %v, want nil`, id, newTestName, err)
	}

	checkExpectedUserData(t, id, tuInstance)
}

func TestUpdateUserStatus(t *testing.T) {
	id, tuInstance := setupUserTest()
	tuInstance.Status = updatedStatus

	err := UpdateUserStatus(id, updatedStatus)
	if err != nil {
		t.Fatalf(`UpdateUserStatus(%s, %s) = %v, want nil`, id, updatedStatus, err)
	}

	checkExpectedUserData(t, id, tuInstance)
}

func setupUserTest() (string, user) {
	id, _ := CreateUser(testUserName)

	tuInstance := testUser
	tuInstance.UID = id

	return id, tuInstance
}

func checkExpectedUserData(t *testing.T, id string, expected user) {
	user, err := GetUser(id)
	if user != expected || err != nil {
		t.Fatalf(`GetUser(UpdateUser) = %v, %v, want %v, nil`, user, err, expected)
	}
}
