package room

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"

	gameclient "Engee-Server/gameClient"
	reg "Engee-Server/gameRegistry"
	sErr "Engee-Server/stockErrors"
	"Engee-Server/testDummy"
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

	if id != "" || !errors.As(err, &sErr.EV_ERR) {
		t.Fatalf(`CreateRoom(EmptyName) = %q, %v, want "", EmptyValueError`, id, err)
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
	if !errors.As(err, &sErr.EV_ERR) {
		t.Fatalf(`GetRoom(EmptyID) = %v, %v, want nil, EmptyValueError`, room, err)
	}
}

func TestGetRoomInvalidID(t *testing.T) {
	setupRoomTest(t)
	room, err := GetRoom(randomID)
	if !errors.As(err, &sErr.MNF_ERR) {
		t.Fatalf(`GetRoom(InvalidID) = %v, %v, want nil, MatchNotFoundError`, room, err)
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
	if !errors.As(err, &sErr.EV_ERR) {
		t.Fatalf(`UpdateRoomName(EmptyName) = %v, want EmptyValueError`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomNameNoChange(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	trInstance.Status = "Created"
	trInstance.Addr = testConURL

	err := UpdateRoomName(id, testRoomName)
	if err != nil {
		t.Fatalf(`UpdateRoomName(NoChange) = %v, want nil`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomNameEmptyID(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	trInstance.Status = "Created"
	trInstance.Addr = testConURL

	err := UpdateRoomName("", altRoomName)
	if !errors.As(err, &sErr.EV_ERR) {
		t.Fatalf(`UpdateRoomName(EmptyID) = %v, want EmptyValueError`, err)
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
	if !errors.As(err, &sErr.EV_ERR) {
		t.Fatalf(`UpdateRoomStatus(EmptyID) = %v, want EmptyValueError`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomStatusInvalidID(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	trInstance.Status = "Created"
	trInstance.Addr = testConURL

	err := UpdateRoomStatus(randomID, updatedRoomStatus)
	if !errors.As(err, &sErr.MNF_ERR) {
		t.Fatalf(`UpdateRoomStatus(InvalidID) = %v, want MatchNotFoundError`, err)
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
	if !errors.As(err, &sErr.EV_ERR) {
		t.Fatalf(`UpdateRoomGameMode(EmptyID) = %v, want EmptyValueError`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
}

func TestUpdateRoomGameModeInvalidID(t *testing.T) {
	id, trInstance := setupRoomTest(t)

	trInstance.Status = "Created"
	trInstance.Addr = testConURL

	err := UpdateRoomGameMode(randomID, testGameMode)
	if !errors.As(err, &sErr.MNF_ERR) {
		t.Fatalf(`UpdateRoomGameMode(InvalidID) = %v, want MatchNotFoundError`, err)
	}

	checkExpectedRoomData(t, id, trInstance)
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
	if !errors.As(err, &sErr.EV_ERR) {
		t.Fatalf(`DeleteRoom(EmptyID) = %v, want EmptyValueError`, err)
	}
}

func TestDeleteInvalidID(t *testing.T) {
	setupActiveRoomTest(t)

	err := DeleteRoom(randomID)
	if !errors.As(err, &sErr.MNF_ERR) {
		t.Fatalf(`DeleteRoom(InvalidID) = %v, want MatchNotFoundError`, err)
	}
}

func TestDeleteDouble(t *testing.T) {
	id, _ := setupActiveRoomTest(t)

	DeleteRoom(id)
	err := DeleteRoom(id)
	if !errors.As(err, &sErr.MNF_ERR) {
		t.Fatalf(`DeleteRoom(Double) = %v, want MatchNotFoundError`, err)
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
