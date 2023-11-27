package room

import "testing"

const testRoomName = "Test Room"
const newRoomName = "New Room"

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
