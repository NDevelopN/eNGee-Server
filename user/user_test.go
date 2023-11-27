package user

import (
	"testing"

	"github.com/google/uuid"
)

var testUser = user{
	UID:    "",
	Name:   testUserName,
	Status: "New",
}

var randomID = uuid.NewString()

const testUserName = "Test User"
const newTestName = "New Name"

const updatedStatus = "Updated"

func TestCreateUser(t *testing.T) {
	id, err := CreateUser(testUserName)
	if id == "" || err != nil {
		t.Fatalf(`CreateUser(%s) = %q, %v, want "uuid", nil`, testUserName, id, err)
	}
}

func TestCreateUserNoName(t *testing.T) {
	id, err := CreateUser("")
	if id != "" || err == nil {
		t.Fatalf(`CreateUser("") = %q, %v, want "", err`, id, err)
	}
}

func TestGetUser(t *testing.T) {
	id, tuInstance := setupUserTest()
	user, err := GetUser(id)
	if user != tuInstance || err != nil {
		t.Fatalf(`GetUser(ValidID) = %v, %v, want obj, nil`, user, err)
	}
}

func TestGetUserEmptyID(t *testing.T) {
	setupUserTest()
	user, err := GetUser("")
	if err == nil {
		t.Fatalf(`GetUser(EmptyID) = %v, %v, want nil, err`, user, err)
	}
}

func TestGetUserInvalidID(t *testing.T) {
	setupUserTest()
	user, err := GetUser(randomID)
	if err == nil {
		t.Fatalf(`GetUser(InvalidID) = %v, %v, want nil, err`, user, err)
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

func TestDeleteUser(t *testing.T) {
	id, _ := setupUserTest()

	err := DeleteUser(id)
	if err != nil {
		t.Fatalf(`DeleteUser(Valid) = %v, want nil`, err)
	}

	confirmUserNotExist(t, id)
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
		t.Fatalf(`GetUser(UpdatedUser) = %v, %v, want %v, nil`, user, err, expected)
	}
}

func confirmUserNotExist(t *testing.T, id string) {
	user, err := GetUser(id)
	if err == nil {
		t.Fatalf(`GetUser(DeletedUser) %v, %v, want nil, err`, user, err)
	}
}
