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
	uid, rid := setupLobbyTest()

	err := JoinUserToRoom(uid, rid)
	if err != nil {
		t.Fatalf(`TestJoinUserToRoom(Valid) = %v, want nil`, err)
	}
}

func TestJoinUserToRoomInvalidUID(t *testing.T) {
	_, rid := setupLobbyTest()

	err := JoinUserToRoom(randomID, rid)
	if err == nil {
		t.Fatalf(`TestJoinUserToRoom(InvalidUID) = %v, want err`, err)
	}
}

func TestJoinUserToRoomInvalidRID(t *testing.T) {
	uid, _ := setupLobbyTest()

	err := JoinUserToRoom(uid, randomID)
	if err == nil {
		t.Fatalf(`TestJoinUserToRoom(InvalidRID) = %v, want err`, err)
	}
}

func TestJoinUserToRoomDouble(t *testing.T) {
	uid, rid := setupLobbyTest()

	JoinUserToRoom(uid, rid)

	err := JoinUserToRoom(uid, rid)
	if err == nil {
		t.Fatalf(`TestJoinUserToRoom(Double) = %v, want err`, err)
	}
}

func setupLobbyTest() (string, string) {
	uid, _ := user.CreateUser(testUserName)
	rid, _ := room.CreateRoom(testRoomName)

	return uid, rid
}
