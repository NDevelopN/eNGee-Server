package lobby

import (
	reg "Engee-Server/gameRegistry"
	"Engee-Server/room"
	"Engee-Server/testDummy"
	"Engee-Server/user"
	"time"

	"encoding/json"
	"os"
	"testing"

	"github.com/google/uuid"
)

var randomID = uuid.NewString()

const testRoomName = "Test Room"
const testConPort = "8091"
const testConURL = "localhost:" + testConPort
const testGameMode = "Test"
const testUserName = "Test User"

var testRoom, _ = json.Marshal(room.Room{
	RID:      "",
	Name:     testRoomName,
	GameMode: testGameMode,
	Status:   "New",
	Addr:     "",
})

const moreUserCount = 3

func TestMain(m *testing.M) {
	setupLobbySuite()
	code := m.Run()
	cleanUpLobbySuite()
	os.Exit(code)
}

func TestJoinUserToRoom(t *testing.T) {
	uid, rid := createUserAndRoom(t)

	err := JoinUserToRoom(uid, rid)
	if err != nil {
		t.Fatalf(`TestJoinUserToRoom(Valid) = %v, want nil`, err)
	}
}

func TestJoinUserToRoomInvalidUID(t *testing.T) {
	_, rid := createUserAndRoom(t)

	err := JoinUserToRoom(randomID, rid)
	if err == nil {
		t.Fatalf(`TestJoinUserToRoom(InvalidUID) = %v, want err`, err)
	}
}

func TestJoinUserToRoomInvalidRID(t *testing.T) {
	uid, _ := createUserAndRoom(t)

	err := JoinUserToRoom(uid, randomID)
	if err == nil {
		t.Fatalf(`TestJoinUserToRoom(InvalidRID) = %v, want err`, err)
	}
}

func TestJoinUserToRoomDouble(t *testing.T) {
	uid, rid := createUserAndRoom(t)

	JoinUserToRoom(uid, rid)

	err := JoinUserToRoom(uid, rid)
	if err == nil {
		t.Fatalf(`TestJoinUserToRoom(Double) = %v, want err`, err)
	}
}

func TestRemoveUserFromRoom(t *testing.T) {
	uid, rid := setupLobbyTest(t)

	err := RemoveUserFromRoom(uid, rid)
	if err != nil {
		t.Fatalf(`TestRemoveUserFromRoom(Valid) = %v, want err`, err)
	}
}

func TestRemoveUserFromRoomInvalidUID(t *testing.T) {
	_, rid := setupLobbyTest(t)

	err := RemoveUserFromRoom(randomID, rid)
	if err == nil {
		t.Fatalf(`TestRemoveUserFromRoom(InvalidUID) = %v, want err`, err)
	}
}

func TestRemoveUserFromRoomInvalidRID(t *testing.T) {
	uid, _ := setupLobbyTest(t)

	err := RemoveUserFromRoom(uid, randomID)
	if err == nil {
		t.Fatalf(`TestRemoveUserFromRoom(InvalidRID) = %v, want err`, err)
	}
}

func TestRemoveUserFromRoomDouble(t *testing.T) {
	uid, rid := setupLobbyTest(t)

	RemoveUserFromRoom(uid, rid)

	err := RemoveUserFromRoom(uid, rid)
	if err == nil {
		t.Fatalf(`TestRemoveUserFromRoom(Double) = %v, want err`, err)
	}

}

func TestGetUsersInRoom(t *testing.T) {
	uid, rid := setupLobbyTest(t)

	users, err := GetUsersInRoom(rid)
	if len(users) != 1 || err != nil {
		t.Fatalf(`TestGetUsersInRoom(Valid) = %v, %v, want [%v], nil`, users, err, uid)
	}
}

func TestGetMultiUsersInRoom(t *testing.T) {
	uid, rid := setupLobbyTest(t)

	var expected = []string{
		uid,
	}

	expected = append(expected, addMoreUsersToLobby(t, rid)...)

	users, err := GetUsersInRoom(rid)
	if len(users) != len(expected) || err != nil {
		t.Fatalf(`TestGetUsersInRoom(Multi) = %v, %v, want %v, nil`, users, err, expected)
	}

	for i, user := range users {
		if user.UID != expected[i] {
			t.Fatalf(`TestGetUsersInRoom(Mult) = %q, want %q`, user.UID, expected[i])
		}
	}
}

func TestGetUsersInRoomInvalidRID(t *testing.T) {
	setupLobbyTest(t)

	users, err := GetUsersInRoom(randomID)
	if len(users) != 0 || err == nil {
		t.Fatalf(`TestGetUsersInRoom(InvalidGID) = %v, %v, want [], err`, users, err)
	}
}

func TestGetUsersInEmptyRoom(t *testing.T) {
	_, rid := createUserAndRoom(t)

	users, err := GetUsersInRoom(rid)
	if len(users) != 0 || err == nil {
		t.Fatalf(`TestGetUsersInRoom(Empty) = %v, %v, want [], err`, users, err)
	}
}

func TestGetUsersInRoomAfterDelete(t *testing.T) {
	uid, rid := setupLobbyTest(t)

	RemoveUserFromRoom(uid, rid)

	users, err := GetUsersInRoom(rid)
	if len(users) != 0 || err == nil {
		t.Fatalf(`TestGetUsersInRoom(AfterUserDelete) = %v, %v, want [], err`, users, err)
	}
}

func TestGetRoomUserCount(t *testing.T) {
	_, rid := setupLobbyTest(t)

	count, err := GetRoomUserCount(rid)
	if count != 1 || err != nil {
		t.Fatalf(`TestGetRoomUserCount(Single) = %d, %v, want 1, nil`, count, err)
	}
}

func TestGetRoomUserCountMulti(t *testing.T) {
	_, rid := setupLobbyTest(t)
	addMoreUsersToLobby(t, rid)

	expected := 1 + moreUserCount
	count, err := GetRoomUserCount(rid)
	if count != expected || err != nil {
		t.Fatalf(`TestGetRoomUserCount(Multi) = %d, %v, want %d, nil`, count, err, expected)
	}
}

func TestGetRoomUserCountInvalidRID(t *testing.T) {
	setupLobbyTest(t)

	count, err := GetRoomUserCount(randomID)
	if count != 0 || err == nil {
		t.Fatalf(`TestGetRoomUserCount(InvalidRID) = %d, %v, want 0, err`, count, err)
	}
}

func TestGetRoomUserCountAfterDelete(t *testing.T) {
	uid, rid := setupLobbyTest(t)

	RemoveUserFromRoom(uid, rid)

	count, err := GetRoomUserCount(rid)
	if count != 0 || err == nil {
		t.Fatalf(`TestGetRoomUserCount(AfterUserDelete) = %d, %v, want 0, err`, count, err)
	}
}

func setupLobbyTest(t *testing.T) (string, string) {
	uid, rid := createUserAndRoom(t)

	JoinUserToRoom(uid, rid)

	t.Cleanup(func() {
		lobbies = make(map[string][]string)
	})

	return uid, rid
}

func createUserAndRoom(t *testing.T) (string, string) {
	uid, err := user.CreateUser(testUserName)
	if err != nil {
		t.Fatalf("Could not create user: %v", err)
	}

	rid, err := room.CreateRoom(testRoom)
	if err != nil {
		t.Fatalf("Could not create room: %v", err)
	}

	t.Cleanup(func() {
		user.DeleteUser(uid)
		room.DeleteRoom(rid)
	})

	return uid, rid
}

func addMoreUsersToLobby(t *testing.T, rid string) []string {
	users := make([]string, 0)
	i := 0
	for i < moreUserCount {
		uid, _ := user.CreateUser(testUserName)
		JoinUserToRoom(uid, rid)
		users = append(users, uid)
		i++
	}

	t.Cleanup(func() {
		for _, uid := range users {
			user.DeleteUser(uid)
		}
	})

	return users
}

func setupLobbySuite() {
	go testDummy.Serve(testConPort)

	reg.RegisterGameMode(testGameMode, testConURL)

	time.Sleep(200 * time.Millisecond)
}

func cleanUpLobbySuite() {

}
