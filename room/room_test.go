package room

import (
	"testing"

	"github.com/google/uuid"
)

var randomID = uuid.NewString()

const testRoomName = "Test Room"
const newRoomName = "New Room"

var testRoom = room{
	RID:     "",
	Name:    testRoomName,
	Type:    "None",
	Status:  "New",
	CurPlrs: 0,
}

func TestCreateRoom(t *testing.T) {
	id, err := CreateRoom(testRoomName)
	if id == "" || err != nil {
		t.Fatalf(`CreateRoom(Valid) = %q, %v, want "uuid", nil`, id, err)
	}
}

func TestCreateUniqueNameRooms(t *testing.T) {
	CreateRoom(testRoomName)
	id, err := CreateRoom(newRoomName)
	if id == "" || err != nil {
		t.Fatalf(`CreateRoom(Unique Name) = %q, %v, want "uuid", nil`, id, err)
	}
}

func TestCreateSameNameRooms(t *testing.T) {
	CreateRoom(testRoomName)
	id, err := CreateRoom(testRoomName)
	if id == "" || err != nil {
		t.Fatalf(`CreateRoom(Same Name) = %q, %v, want "uuid", nil`, id, err)
	}
}

func TestCreateRoomNoName(t *testing.T) {
	id, err := CreateRoom("")
	if id != "" || err == nil {
		t.Fatalf(`CreateRoom(EmptyName) = %q, %v, want "", nil`, id, err)
	}
}

func TestGetRoom(t *testing.T) {
	id, trInstance := setupRoomTest()

	room, err := GetRoom(id)
	if room != trInstance || err != nil {
		t.Fatalf(`GetRoom(ValidID) = %v, %v, want %v, nil`, room, err, trInstance)
	}
}

func TestGetRoomEmptyID(t *testing.T) {
	setupRoomTest()
	room, err := GetRoom("")
	if err == nil {
		t.Fatalf(`GetRoom(EmptyID) = %v, %v, want nil, err`, room, err)
	}
}

func TestGetUserInvalidID(t *testing.T) {
	setupRoomTest()
	room, err := GetRoom(randomID)
	if err == nil {
		t.Fatalf(`GetUser(InvalidID) = %v, %v, want nil, err`, room, err)
	}
}

func setupRoomTest() (string, room) {
	id, _ := CreateRoom(testRoomName)

	tuInstance := testRoom
	tuInstance.RID = id

	return id, tuInstance
}
