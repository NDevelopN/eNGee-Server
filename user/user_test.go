package user

import (
	"testing"

	"github.com/google/uuid"
)

var testUser = User{
	UID:    "",
	Name:   testUserName,
	Status: "New",
}

var randomID = uuid.NewString()

const testUserName = "Test User"
const newUserName = "New Name"

const updatedUserStatus = "Updated"

func TestCreateUser(t *testing.T) {
	id, err := CreateUser(testUserName)
	if id == "" || err != nil {
		t.Fatalf(`CreateUser(Valid) = %q, %v, want "uuid", nil`, id, err)
	}

	t.Cleanup(cleanAfterTest)
}

func TestCreateUniqueNameUsers(t *testing.T) {
	CreateUser(testUserName)
	id, err := CreateUser(newUserName)
	if id == "" || err != nil {
		t.Fatalf(`CreateUser(Unique Name) = %q, %v, want "uuid", nil`, id, err)
	}

	t.Cleanup(cleanAfterTest)
}

func TestCreateSameNameUsers(t *testing.T) {
	CreateUser(testUserName)
	id, err := CreateUser(testUserName)
	if id == "" || err != nil {
		t.Fatalf(`CreateUser(Same Name) = %q, %v, want "uuid", nil`, id, err)
	}

	t.Cleanup(cleanAfterTest)
}

func TestCreateUserEmptyName(t *testing.T) {
	id, err := CreateUser("")
	if id != "" || err == nil {
		t.Fatalf(`CreateUser(EmptyName) = %q, %v, want "", err`, id, err)
	}

	t.Cleanup(cleanAfterTest)
}

func TestGetUser(t *testing.T) {
	id, tuInstance := setupUserTest(t)
	user, err := GetUser(id)
	if user != tuInstance || err != nil {
		t.Fatalf(`GetUser(ValidID) = %v, %v, want obj, nil`, user, err)
	}
}

func TestGetUserEmptyID(t *testing.T) {
	setupUserTest(t)
	user, err := GetUser("")
	if err == nil {
		t.Fatalf(`GetUser(EmptyID) = %v, %v, want nil, err`, user, err)
	}
}

func TestGetUserInvalidID(t *testing.T) {
	setupUserTest(t)
	user, err := GetUser(randomID)
	if err == nil {
		t.Fatalf(`GetUser(InvalidID) = %v, %v, want nil, err`, user, err)
	}
}

func TestUpdateUserName(t *testing.T) {
	id, tuInstance := setupUserTest(t)

	tuInstance.Name = newUserName

	err := UpdateUserName(id, newUserName)
	if err != nil {
		t.Fatalf(`UpdateUserName(Valid) = %v, want nil`, err)
	}

	checkExpectedUserData(t, id, tuInstance)
}

func TestUpdateUserNameEmptyName(t *testing.T) {
	id, tuInstance := setupUserTest(t)

	err := UpdateUserName(id, "")
	if err == nil {
		t.Fatalf(`UpdateUserName(EmptyName) = %v, want err`, err)
	}

	checkExpectedUserData(t, id, tuInstance)
}

func TestUpdateUserNameNoChange(t *testing.T) {
	id, tuInstance := setupUserTest(t)

	err := UpdateUserName(id, testUserName)
	if err != nil {
		t.Fatalf(`UpdateUserName(NoChange) = %v, want err`, err)
	}

	checkExpectedUserData(t, id, tuInstance)
}
func TestUpdateUserNameEmptyID(t *testing.T) {
	id, tuInstance := setupUserTest(t)

	err := UpdateUserName("", newUserName)
	if err == nil {
		t.Fatalf(`UpdateUserName(EmptyID) = %v, want err`, err)
	}

	checkExpectedUserData(t, id, tuInstance)
}

func TestUpdateUserStatusEmptyStatus(t *testing.T) {
	id, tuInstance := setupUserTest(t)

	err := UpdateUserStatus(id, "")
	if err == nil {
		t.Fatalf(`UpdateUserStatus(EmptyStatus) = %v, want err`, err)
	}

	checkExpectedUserData(t, id, tuInstance)
}
func TestUpdateUserStatusNoChange(t *testing.T) {
	id, tuInstance := setupUserTest(t)

	err := UpdateUserStatus(id, testUser.Status)
	if err != nil {
		t.Fatalf(`UpdateUserStatus(NoChange) = %v, want err`, err)
	}

	checkExpectedUserData(t, id, tuInstance)
}
func TestUpdateUserNameInvalidID(t *testing.T) {
	id, tuInstance := setupUserTest(t)

	err := UpdateUserName(randomID, newUserName)
	if err == nil {
		t.Fatalf(`UpdateUserName(InvalidID) = %v, want err`, err)
	}

	checkExpectedUserData(t, id, tuInstance)
}

func TestUpdateUserStatus(t *testing.T) {
	id, tuInstance := setupUserTest(t)
	tuInstance.Status = updatedUserStatus

	err := UpdateUserStatus(id, updatedUserStatus)
	if err != nil {
		t.Fatalf(`UpdateUserStatus(Valid) = %v, want nil`, err)
	}

	checkExpectedUserData(t, id, tuInstance)
}

func TestUpdateUserStatusEmptyID(t *testing.T) {
	id, tuInstance := setupUserTest(t)

	err := UpdateUserName("", newUserName)
	if err == nil {
		t.Fatalf(`UpdateUserStatus(EmptyID) = %v, want err`, err)
	}

	checkExpectedUserData(t, id, tuInstance)
}

func TestUpdateUserStatusInvalidID(t *testing.T) {
	id, tuInstance := setupUserTest(t)

	err := UpdateUserStatus(randomID, updatedUserStatus)
	if err == nil {
		t.Fatalf(`UpdateUserStatus(InvalidID) = %v, want err`, err)
	}

	checkExpectedUserData(t, id, tuInstance)
}

func TestDeleteUser(t *testing.T) {
	id, _ := setupUserTest(t)

	err := DeleteUser(id)
	if err != nil {
		t.Fatalf(`DeleteUser(Valid) = %v, want nil`, err)
	}

	confirmUserNotExist(t, id)
}

func TestDeleteEmptyID(t *testing.T) {
	setupUserTest(t)

	err := DeleteUser("")
	if err == nil {
		t.Fatalf(`DeleteUser(EmptyID) = %v, want err`, err)
	}
}

func TestDeleteInvalidID(t *testing.T) {
	setupUserTest(t)

	err := DeleteUser(randomID)
	if err == nil {
		t.Fatalf(`DeleteUser(InvalidID) = %v, want err`, err)
	}
}

func TestDeleteDouble(t *testing.T) {
	id, _ := setupUserTest(t)

	DeleteUser(id)
	err := DeleteUser(id)
	if err == nil {
		t.Fatalf(`DeleteUser(Double) = %v, want err`, err)
	}
}

func setupUserTest(t *testing.T) (string, User) {
	id, _ := CreateUser(testUserName)

	tuInstance := testUser
	tuInstance.UID = id

	t.Cleanup(cleanAfterTest)

	return id, tuInstance
}

func checkExpectedUserData(t *testing.T, id string, expected User) {
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

func cleanAfterTest() {
	users = make(map[string]User)
}
