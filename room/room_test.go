package room

import (
	"testing"

	"github.com/google/uuid"
)

var randomID = uuid.NewString()

const testRoomName = "Test Room"
const newRoomName = "New Room"

const updatedRoomStatus = "Updated"
const updatedRoomType = "New Type"

var testRoom = room{
	RID:    "",
	Name:   testRoomName,
	Type:   "None",
	Status: "New",
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

func TestGetRoomInvalidID(t *testing.T) {
	setupRoomTest()
	room, err := GetRoom(randomID)
	if err == nil {
		t.Fatalf(`GetRoom(InvalidID) = %v, %v, want nil, err`, room, err)
	}
}

func TestGetRooms(t *testing.T) {
	fID, fRoom := setupRoomTest()
	sID, sRoom := setupAddRoomTest()
	expected := map[string]room{fID: fRoom, sID: sRoom}

	rooms, err := GetRooms()
	if len(rooms) != 2 || err != nil {
		t.Fatalf(`GetRooms(Valid) = %v, %v, want %v, nil`, rooms, err, expected)
	}

	for id, r := range rooms {
		if r != expected[id] {
			t.Fatalf(`GetRooms(Valid) = %v, want %v`, r, expected[id])
		}
	}
}

func TestGetRoomsEmpty(t *testing.T) {
	rooms, err := GetRooms()
	if len(rooms) != 0 || err == nil {
		t.Fatalf(`GetRooms(Empty) = %v, %v, want [], err`, rooms, err)
	}
}
func TestUpdateRoomName(t *testing.T) {
	id, trInstance := setupRoomTest()

	trInstance.Name = newRoomName

	err := UpdateRoomName(id, newRoomName)
	if err != nil {
		t.Fatalf(`UpdateRoomName(Valid) = %v, want nil`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomNameEmptyName(t *testing.T) {
	id, trInstance := setupRoomTest()

	err := UpdateRoomName(id, "")
	if err == nil {
		t.Fatalf(`UpdateRoomName(EmptyName) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomNameNoChange(t *testing.T) {
	id, trInstance := setupRoomTest()

	err := UpdateRoomName(id, testRoomName)
	if err != nil {
		t.Fatalf(`UpdateRoomName(NoChange) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}
func TestUpdateRoomNameEmptyID(t *testing.T) {
	id, trInstance := setupRoomTest()

	err := UpdateRoomName("", newRoomName)
	if err == nil {
		t.Fatalf(`UpdateRoomName(EmptyID) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomStatus(t *testing.T) {
	id, trInstance := setupRoomTest()
	trInstance.Status = updatedRoomStatus

	err := UpdateRoomStatus(id, updatedRoomStatus)
	if err != nil {
		t.Fatalf(`UpdateRoomStatus(%s, %s) = %v, want nil`, id, updatedRoomStatus, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomStatusEmptyID(t *testing.T) {
	id, trInstance := setupRoomTest()

	err := UpdateRoomName("", newRoomName)
	if err == nil {
		t.Fatalf(`UpdateRoomStatus(EmptyID) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomStatusInvalidID(t *testing.T) {
	id, trInstance := setupRoomTest()

	err := UpdateRoomStatus(randomID, updatedRoomStatus)
	if err == nil {
		t.Fatalf(`UpdateRoomStatus(InvalidID) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomType(t *testing.T) {
	id, trInstance := setupRoomTest()
	trInstance.Type = updatedRoomType

	err := UpdateRoomType(id, updatedRoomType)
	if err != nil {
		t.Fatalf(`UpdateRoomType(%s, %s) = %v, want nil`, id, updatedRoomType, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomTypeEmptyID(t *testing.T) {
	id, trInstance := setupRoomTest()

	err := UpdateRoomName("", newRoomName)
	if err == nil {
		t.Fatalf(`UpdateRoomType(EmptyID) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomTypeInvalidID(t *testing.T) {
	id, trInstance := setupRoomTest()

	err := UpdateRoomType(randomID, updatedRoomType)
	if err == nil {
		t.Fatalf(`UpdateRoomType(InvalidID) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func setupRoomTest() (string, room) {
	id, _ := CreateRoom(testRoomName)

	trInstance := testRoom
	trInstance.RID = id

	return id, trInstance
}

func TestDeleteRoom(t *testing.T) {
	id, _ := setupRoomTest()

	err := DeleteRoom(id)
	if err != nil {
		t.Fatalf(`DeleteRoom(Valid) = %v, want nil`, err)
	}

	confirmRoomNotExist(t, id)
}

func TestDeleteEmptyID(t *testing.T) {
	setupRoomTest()

	err := DeleteRoom("")
	if err == nil {
		t.Fatalf(`DeleteRoom(EmptyID) = %v, want err`, err)
	}
}

func TestDeleteInvalidID(t *testing.T) {
	setupRoomTest()

	err := DeleteRoom(randomID)
	if err == nil {
		t.Fatalf(`DeleteRoom(InvalidID) = %v, want err`, err)
	}
}

func TestDeleteDouble(t *testing.T) {
	id, _ := setupRoomTest()

	DeleteRoom(id)
	err := DeleteRoom(id)
	if err == nil {
		t.Fatalf(`DeleteRoom(Double) = %v, want err`, err)
	}
}

func setupAddRoomTest() (string, room) {
	id, _ := CreateRoom(newRoomName)

	trInstance := testRoom
	trInstance.Name = newRoomName
	trInstance.RID = id

	return id, trInstance
}

func checkExpectedRoomData(t *testing.T, id string, expected room) {
	room, err := GetRoom(id)
	if room != expected || err != nil {
		t.Fatalf(`GetRoom(UpdatedRoom) = %v, %v, want %v, nil`, room, err, expected)
	}
}

func confirmRoomNotExist(t *testing.T, id string) {
	room, err := GetRoom(id)
	if err == nil {
		t.Fatalf(`GetRoom(DeletedRoom) %v, %v, want nil, err`, room, err)
	}
}
