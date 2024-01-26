package room

import (
	gameclient "Engee-Server/gameClient"
	reg "Engee-Server/gameRegistry"
	"Engee-Server/testDummy"
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

var randomID = uuid.NewString()

const testRoomName = "Test-Room"
const altRoomName = "Alt-Room"

const updatedRoomStatus = "Updated"

const testGameMode = "Test"
const altGameMode = "Alt"

const testConPort = "8091"
const altConPort = "8092"
const testConURL = "http://localhost:" + testConPort
const altConURL = "http://localhost:" + altConPort

var testRoom = Room{
	RID:      "",
	Name:     testRoomName,
	GameMode: testGameMode,
	Status:   "New",
	Addr:     "",
}

var testRoomJSON, _ = json.Marshal(testRoom)

var altRoom = Room{
	RID:      "",
	Name:     altRoomName,
	GameMode: altGameMode,
	Status:   "New",
	Addr:     "",
}

var altRoomJSON, _ = json.Marshal(altRoom)

func TestMain(m *testing.M) {
	setupRoomSuite()
	code := m.Run()
	cleanUpAfterSuite()
	os.Exit(code)
}

func TestCreateRoom(t *testing.T) {
	id, err := CreateRoom(testRoomJSON)
	if id == "" || err != nil {
		t.Fatalf(`CreateRoom(Valid) = %q, %v, want "uuid", nil`, id, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateUniqueNameRooms(t *testing.T) {
	CreateRoom(testRoomJSON)

	id, err := CreateRoom(testRoomJSON)

	if id == "" || err != nil {
		t.Fatalf(`CreateRoom(Unique Name) = %q, %v, want "uuid", nil`, id, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateSameNameRooms(t *testing.T) {
	CreateRoom(testRoomJSON)
	id, err := CreateRoom(testRoomJSON)
	if id == "" || err != nil {
		t.Fatalf(`CreateRoom(Same Name) = %q, %v, want "uuid", nil`, id, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestCreateRoomNoName(t *testing.T) {
	namelessRoom, _ := json.Marshal(Room{
		Name:     "",
		GameMode: "None",
	})

	id, err := CreateRoom(namelessRoom)

	if id != "" || err == nil {
		t.Fatalf(`CreateRoom(EmptyName) = %q, %v, want "", nil`, id, err)
	}

	t.Cleanup(cleanUpAfterTest)
}

func TestGetRoom(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	trInstance.Status = "Created"
	trInstance.Addr = testConURL

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
	_, fRoom := setupRoomTest(t)
	fRoom.Status = "Created"
	fRoom.Addr = testConURL

	_, sRoom := setupAltRoomTest()
	sRoom.Status = "Created"
	sRoom.Addr = altConURL

	expected := []Room{fRoom, sRoom}

	rooms := GetRooms()
	if len(rooms) != 2 {
		t.Fatalf(`GetRooms(Valid) = %v, want %v`, rooms, expected)
	}

	unmatched := len(expected)

	if len(rooms) == unmatched {
		for _, r := range rooms {
			log.Printf("R: %s", r.RID)
			for _, e := range expected {
				log.Printf("E: %s", e.RID)
				if r == e {
					log.Printf("Matching %s", r.RID)
					unmatched--
				}
			}
		}
	}

	if unmatched != 0 {
		t.Fatalf(`GetRooms(Valid) = %v, want %v`, rooms, expected)
	}
}

func TestGetRoomsEmpty(t *testing.T) {
	rooms := GetRooms()
	if len(rooms) != 0 {
		t.Fatalf(`GetRooms(Empty) = %v, want []`, rooms)
	}
}

func TestUpdateRoomName(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	trInstance.Name = altRoomName
	trInstance.Status = "Created"
	trInstance.Addr = testConURL

	err := UpdateRoomName(id, altRoomName)
	if err != nil {
		t.Fatalf(`UpdateRoomName(Valid) = %v, want nil`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomNameEmptyName(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	trInstance.Status = "Created"
	trInstance.Addr = testConURL

	err := UpdateRoomName(id, "")
	if err == nil {
		t.Fatalf(`UpdateRoomName(EmptyName) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomNameNoChange(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	trInstance.Status = "Created"
	trInstance.Addr = testConURL

	err := UpdateRoomName(id, testRoomName)
	if err != nil {
		t.Fatalf(`UpdateRoomName(NoChange) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomNameEmptyID(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	trInstance.Status = "Created"
	trInstance.Addr = testConURL

	err := UpdateRoomName("", altRoomName)
	if err == nil {
		t.Fatalf(`UpdateRoomName(EmptyID) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomStatus(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	trInstance.Status = updatedRoomStatus
	trInstance.Addr = testConURL

	err := UpdateRoomStatus(id, updatedRoomStatus)
	if err != nil {
		t.Fatalf(`UpdateRoomStatus(Valid) = %v, want nil`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomStatusEmptyID(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	trInstance.Status = "Created"
	trInstance.Addr = testConURL

	err := UpdateRoomName("", altRoomName)
	if err == nil {
		t.Fatalf(`UpdateRoomStatus(EmptyID) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomStatusInvalidID(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	trInstance.Status = "Created"
	trInstance.Addr = testConURL

	err := UpdateRoomStatus(randomID, updatedRoomStatus)
	if err == nil {
		t.Fatalf(`UpdateRoomStatus(InvalidID) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomGameMode(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	trInstance.Status = "Created"
	trInstance.GameMode = altGameMode
	trInstance.Addr = altConURL

	err := UpdateRoomGameMode(id, altGameMode)
	if err != nil {
		t.Fatalf(`UpdateRoomGameMode(Valid) = %v, want nil`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomGameModeEmptyID(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	trInstance.Status = "Created"
	trInstance.Addr = testConURL

	err := UpdateRoomName("", altRoomName)
	if err == nil {
		t.Fatalf(`UpdateRoomGameMode(EmptyID) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomGameModeInvalidID(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	trInstance.Status = "Created"
	trInstance.Addr = testConURL

	err := UpdateRoomGameMode(randomID, testGameMode)
	if err == nil {
		t.Fatalf(`UpdateRoomGameMode(InvalidID) = %v, want err`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

/** TODO: Given auto initialization, is there any need for this function ? */
/**
func TestInitializeRoomGame(t *testing.T) {
	id, _ := setupRoomTest(t)
	UpdateRoomGameMode(id, testGameMode)

	err := InitializeRoomGame(id)
	if err != nil {
		t.Fatalf(`CreateRoomGameInstance(Valid) = %v, want nil`, err)
	}
}

func TestInitializeRoomGameDouble(t *testing.T) {
	id, _ := setupRoomTest(t)
	UpdateRoomGameMode(id, testGameMode)

	InitializeRoomGame(id)
	err := InitializeRoomGame(id)
	if err == nil {
		t.Fatalf(`CreateRoomGameInstance(Double) = %v, want err`, err)
	}
}
func TestInitializeRoomGameInvalidRID(t *testing.T) {
	id, _ := setupRoomTest(t)
	UpdateRoomGameMode(id, testGameMode)

	err := InitializeRoomGame(randomID)
	if err == nil {
		t.Fatalf(`CreateRoomGameInstance(Invalid RID) = %v, want err`, err)
	}
}
func TestInitializeRoomGameModeNotSet(t *testing.T) {
	id, _ := setupRoomTest(t)

	err := InitializeRoomGame(id)
	if err == nil {
		t.Fatalf(`CreateRoomGameInstance(Room GameMode Not Set) = %v, want err`, err)
	}
}

func TestInitializeRoomGameDeletedRoom(t *testing.T) {
	id, _ := setupRoomTest(t)
	UpdateRoomGameMode(id, testGameMode)

	err := DeleteRoom(id)
	if err != nil {
		t.Fatalf(`Failed to delete: %v`, err)
	}
	err = InitializeRoomGame(id)
	if err == nil {
		t.Fatalf(`CreateRoomGameInstance(Deleted Room) = %v, want err`, err)
	}
}
*/

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
	id, _ := CreateRoom(testRoomJSON)

	trInstance := testRoom
	trInstance.RID = id

	t.Cleanup(cleanUpAfterTest)

	return id, trInstance
}

func setupAltRoomTest() (string, Room) {
	id, _ := CreateRoom(altRoomJSON)

	trInstance := altRoom
	trInstance.RID = id

	return id, trInstance
}

func setupActiveRoomTest(t *testing.T) (string, Room) {
	id, _ := setupRoomTest(t)

	UpdateRoomGameMode(id, testGameMode)

	trInstance, _ := GetRoom(id)

	gameclient.CreateGameInstance(id, trInstance.Addr)

	return id, trInstance
}

func setupRoomSuite() {
	go testDummy.Serve(testConPort)
	go testDummy.Serve(altConPort)

	reg.RegisterGameMode(testGameMode, testConURL)
	reg.RegisterGameMode(altGameMode, altConURL)

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
