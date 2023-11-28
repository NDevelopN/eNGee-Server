package lobby

import (
	"Engee-Server/room"
	"Engee-Server/user"
	"testing"

	"github.com/google/uuid"
)

var randomID = uuid.NewString()

const testRoomName = "Test Room"
const testUserName = "Test User"

func TestJoinUserToRoom(t *testing.T) {
	uid, rid := createUserAndRoom()

	err := JoinUserToRoom(uid, rid)
	if err != nil {
		t.Fatalf(`TestJoinUserToRoom(Valid) = %v, want nil`, err)
	}
}

func TestJoinUserToRoomInvalidUID(t *testing.T) {
	_, rid := createUserAndRoom()

	err := JoinUserToRoom(randomID, rid)
	if err == nil {
		t.Fatalf(`TestJoinUserToRoom(InvalidUID) = %v, want err`, err)
	}
}

func TestJoinUserToRoomInvalidRID(t *testing.T) {
	uid, _ := createUserAndRoom()

	err := JoinUserToRoom(uid, randomID)
	if err == nil {
		t.Fatalf(`TestJoinUserToRoom(InvalidRID) = %v, want err`, err)
	}
}

func TestJoinUserToRoomDouble(t *testing.T) {
	uid, rid := createUserAndRoom()

	JoinUserToRoom(uid, rid)

	err := JoinUserToRoom(uid, rid)
	if err == nil {
		t.Fatalf(`TestJoinUserToRoom(Double) = %v, want err`, err)
	}
}

func TestRemoveUserFromRoom(t *testing.T) {
	uid, rid := setupLobbyTest()

	err := RemoveUserFromRoom(uid, rid)
	if err != nil {
		t.Fatalf(`TestRemoveUserFromRoom(Valid) = %v, want err`, err)
	}
}

func TestRemoveUserFromRoomInvalidUID(t *testing.T) {
	_, rid := setupLobbyTest()

	err := RemoveUserFromRoom(randomID, rid)
	if err == nil {
		t.Fatalf(`TestRemoveUserFromRoom(InvalidUID) = %v, want err`, err)
	}
}

func TestRemoveUserFromRoomInvalidRID(t *testing.T) {
	uid, _ := setupLobbyTest()

	err := RemoveUserFromRoom(uid, randomID)
	if err == nil {
		t.Fatalf(`TestRemoveUserFromRoom(InvalidRID) = %v, want err`, err)
	}
}

func TestRemoveUserFromRoomDouble(t *testing.T) {
	uid, rid := setupLobbyTest()

	RemoveUserFromRoom(uid, rid)

	err := RemoveUserFromRoom(uid, rid)
	if err == nil {
		t.Fatalf(`TestRemoveUserFromRoom(Double) = %v, want err`, err)
	}
}

func setupLobbyTest() (string, string) {
	uid, rid := createUserAndRoom()

	JoinUserToRoom(uid, rid)

	return uid, rid
}

func createUserAndRoom() (string, string) {
	uid, _ := user.CreateUser(testUserName)
	rid, _ := room.CreateRoom(testRoomName)

	return uid, rid
}
