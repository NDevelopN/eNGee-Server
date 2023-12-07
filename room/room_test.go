package room

import (
	gameclient "Engee-Server/gameClient"
	reg "Engee-Server/gameRegistry"
	"os"
	"testing"
	"time"

	"Engee-Server/gamedummy"

	"github.com/google/uuid"
)

var randomID = uuid.NewString()

const testRoomName = "Test Room"
const newRoomName = "New Room"

const updatedRoomStatus = "Updated"

const testRoomType = "Test"
const altRoomType = "Alt"

const testConPort = "8091"
const altConPort = "8092"
const testConURL = "http://localhost:" + testConPort
const altConURL = "http://localhost:" + altConPort

var testRoom = Room{
	RID:    "",
	Name:   testRoomName,
	Type:   "None",
	Status: "New",
	Addr:   "",
}

func TestMain(m *testing.M) {
	setupRoomSuite()
	code := m.Run()
	cleanUpAfterSuite()
	os.Exit(code)
}

func TestCreateRoom(t *testing.T) {
	id, err := CreateRoom(testRoomName)
	if id == "" || err != nil {
		t.Fatalf(`CreateRoom(Valid) = %q, %v, want "uuid", nil`, id, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateUniqueNameRooms(t *testing.T) {
	CreateRoom(testRoomName)

	id, err := CreateRoom(newRoomName)

	if id == "" || err != nil {
		t.Fatalf(`CreateRoom(Unique Name) = %q, %v, want "uuid", nil`, id, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateSameNameRooms(t *testing.T) {
	CreateRoom(testRoomName)
	id, err := CreateRoom(testRoomName)
	if id == "" || err != nil {
		t.Fatalf(`CreateRoom(Same Name) = %q, %v, want "uuid", nil`, id, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateRoomNoName(t *testing.T) {
	id, err := CreateRoom("")
	if id != "" || err == nil {
		t.Fatalf(`CreateRoom(EmptyName) = %q, %v, want "", nil`, id, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestGetRoom(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	room, err := GetRoom(id)
	if room != trInstance || err != nil {
		t.Fatalf(`GetRoom(ValidID) = %v, %v, want %v, nil`, room, err, trInstance)
	}
}

func TestGetRoomEmptyID(t *testing.T) {
	setupRoomTest(t)
	room, err := GetRoom("")
	if err == nil {
		t.Fatalf(`GetRoom(EmptyID) = %v, %v, want nil, err`, room, err)
	}
}

func TestGetRoomInvalidID(t *testing.T) {
	setupRoomTest(t)
	room, err := GetRoom(randomID)
	if err == nil {
		t.Fatalf(`GetRoom(InvalidID) = %v, %v, want nil, err`, room, err)
	}
}

func TestGetRooms(t *testing.T) {
	fID, fRoom := setupRoomTest(t)
	sID, sRoom := setupAddRoomTest()
	expected := map[string]Room{fID: fRoom, sID: sRoom}

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
	id, trInstance := setupRoomTest(t)

	trInstance.Name = newRoomName

	err := UpdateRoomName(id, newRoomName)
	if err != nil {
		t.Fatalf(`UpdateRoomName(Valid) = %v, want nil`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomNameEmptyName(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	err := UpdateRoomName(id, "")
	if err == nil {
		t.Fatalf(`UpdateRoomName(EmptyName) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomNameNoChange(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	err := UpdateRoomName(id, testRoomName)
	if err != nil {
		t.Fatalf(`UpdateRoomName(NoChange) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomNameEmptyID(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	err := UpdateRoomName("", newRoomName)
	if err == nil {
		t.Fatalf(`UpdateRoomName(EmptyID) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomStatus(t *testing.T) {
	id, trInstance := setupRoomTest(t)
	trInstance.Status = updatedRoomStatus

	err := UpdateRoomStatus(id, updatedRoomStatus)
	if err != nil {
		t.Fatalf(`UpdateRoomStatus(Valid) = %v, want nil`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomStatusEmptyID(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	err := UpdateRoomName("", newRoomName)
	if err == nil {
		t.Fatalf(`UpdateRoomStatus(EmptyID) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomStatusInvalidID(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	err := UpdateRoomStatus(randomID, updatedRoomStatus)
	if err == nil {
		t.Fatalf(`UpdateRoomStatus(InvalidID) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomType(t *testing.T) {
	id, trInstance := setupRoomTest(t)
	trInstance.Type = testRoomType
	trInstance.Addr = testConURL

	err := UpdateRoomType(id, testRoomType)
	if err != nil {
		t.Fatalf(`UpdateRoomType(Valid) = %v, want nil`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomTypeEmptyID(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	err := UpdateRoomName("", newRoomName)
	if err == nil {
		t.Fatalf(`UpdateRoomType(EmptyID) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomTypeInvalidID(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	err := UpdateRoomType(randomID, testRoomType)
	if err == nil {
		t.Fatalf(`UpdateRoomType(InvalidID) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestCreateRoomInstance(t *testing.T) {
	id, _ := setupRoomTest(t)
	UpdateRoomType(id, testRoomType)

	err := CreateRoomInstance(id)
	if err != nil {
		t.Fatalf(`CreateRoomGameInstance(Valid) = %v, want nil`, err)
	}
}

func TestCreateRoomInstanceDouble(t *testing.T) {
	id, _ := setupRoomTest(t)
	UpdateRoomType(id, testRoomType)

	CreateRoomInstance(id)
	err := CreateRoomInstance(id)
	if err == nil {
		t.Fatalf(`CreateRoomGameInstance(Double) = %v, want err`, err)
	}
}
func TestCreateRoomInstanceInvalidRID(t *testing.T) {
	id, _ := setupRoomTest(t)
	UpdateRoomType(id, testRoomType)

	err := CreateRoomInstance(randomID)
	if err == nil {
		t.Fatalf(`CreateRoomGameInstance(Invalid RID) = %v, want err`, err)
	}
}
func TestCreateRoomInstanceTypeNotSet(t *testing.T) {
	id, _ := setupRoomTest(t)

	err := CreateRoomInstance(id)
	if err == nil {
		t.Fatalf(`CreateRoomGameInstance(Room Type Not Set) = %v, want err`, err)
	}
}

func TestCreateRoomInstanceDeletedRoom(t *testing.T) {
	id, _ := setupRoomTest(t)
	UpdateRoomType(id, testRoomType)

	err := DeleteRoom(id)
	if err != nil {
		t.Fatalf(`Failed to delete: %v`, err)
	}
	err = CreateRoomInstance(id)
	if err == nil {
		t.Fatalf(`CreateRoomGameInstance(Deleted Room) = %v, want err`, err)
	}
}

func TestDeleteRoom(t *testing.T) {
	id, _ := setupActiveRoomTest(t)

	err := DeleteRoom(id)
	if err != nil {
		t.Fatalf(`DeleteRoom(Valid) = %v, want nil`, err)
	}

	confirmRoomNotExist(t, id)
}

func TestDeleteEmptyID(t *testing.T) {
	setupActiveRoomTest(t)

	err := DeleteRoom("")
	if err == nil {
		t.Fatalf(`DeleteRoom(EmptyID) = %v, want err`, err)
	}
}

func TestDeleteInvalidID(t *testing.T) {
	setupActiveRoomTest(t)

	err := DeleteRoom(randomID)
	if err == nil {
		t.Fatalf(`DeleteRoom(InvalidID) = %v, want err`, err)
	}
}

func TestDeleteDouble(t *testing.T) {
	id, _ := setupActiveRoomTest(t)

	DeleteRoom(id)
	err := DeleteRoom(id)
	if err == nil {
		t.Fatalf(`DeleteRoom(Double) = %v, want err`, err)
	}
}

func setupRoomTest(t *testing.T) (string, Room) {
	id, _ := CreateRoom(testRoomName)

	trInstance := testRoom
	trInstance.RID = id

	t.Cleanup(cleanUpAfterTest)

	return id, trInstance
}

func setupAddRoomTest() (string, Room) {
	id, _ := CreateRoom(newRoomName)

	trInstance := testRoom
	trInstance.Name = newRoomName
	trInstance.RID = id

	return id, trInstance
}

func setupActiveRoomTest(t *testing.T) (string, Room) {
	id, _ := setupRoomTest(t)

	UpdateRoomType(id, testRoomType)

	trInstance, _ := GetRoom(id)

	gameclient.CreateGame(id, trInstance.Addr)

	return id, trInstance
}

func setupRoomSuite() {
	go gamedummy.Start(testConPort)
	go gamedummy.Start(altConPort)

	reg.RegisterGameType(testRoomType, testConURL)
	reg.RegisterGameType(altRoomType, altConURL)

	time.Sleep(200 * time.Millisecond)
}

func checkExpectedRoomData(t *testing.T, id string, expected Room) {
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

func cleanUpAfterTest() {
	rooms = make(map[string]Room)
}

func cleanUpAfterSuite() {

}
