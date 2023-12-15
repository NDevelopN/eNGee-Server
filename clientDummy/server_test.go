package clientdummy

import (
	gameclient "Engee-Server/gameClient"
	registry "Engee-Server/gameRegistry"
	"Engee-Server/lobby"
	"Engee-Server/room"
	"Engee-Server/server"
	"Engee-Server/user"
	"encoding/json"
	"fmt"

	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

const testMainPort = "8080"
const testMainURL = "http://localhost:" + testMainPort

const testConPort = "8081"
const testConURL = "http://localhost:" + testConPort

const testUser = "Test User"
const altUser = "Alt User"

const testRoom = "Test Room"
const altRoom = "Alt Room"

const newRules = "New Rules"

const testType = "Test"

func TestMain(m *testing.M) {
	setupClientSuite()
	code := m.Run()
	cleanUpAfterSuite()
	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	resp, err := sendRequest(testMainURL+"/users", http.MethodPost, []byte(testUser))
	if err != nil {
		t.Fatalf(`TestCreateUser() = %v, want nil`, err)
	}

	uid, err := uuid.Parse(resp)
	if err != nil {
		t.Fatalf(`TestCreateUser() = %q, %v want uuid, nil`, uid, err)
	}

	_, err = user.GetUser(resp)
	if err != nil {
		t.Fatalf(`TestCreateUser() = %v, want nil`, err)
	}
}

func TestCreateRoom(t *testing.T) {
	resp, err := sendRequest(testMainURL+"/rooms", http.MethodPost, []byte(testRoom))
	if err != nil {
		t.Fatalf(`TestCreateRoom() = %v, want nil`, err)
	}

	rid, err := uuid.Parse(resp)
	if err != nil {
		t.Fatalf(`testCreateRoom() = %q, %v want uuid, nil`, rid, err)
	}

	_, err = room.GetRoom(resp)
	if err != nil {
		t.Fatalf(`TestCreateRoom() = %v, want nil`, err)
	}
}

func TestGetRooms(t *testing.T) {
	rid := setupRoom(t)

	resp, err := sendRequest(testMainURL+"/rooms", http.MethodGet, []byte{})
	if err != nil {
		t.Fatalf(`TestGetRooms() = %v, want nil`, err)
	}

	var rooms map[string]room.Room
	err = json.Unmarshal([]byte(resp), &rooms)

	if err != nil {
		t.Fatalf(`TestGetRooms() = %v, %v, want roomStruct, nil`, rooms, err)
	}

	_, ok := rooms[rid]
	if !ok {
		t.Fatalf(`TestGetRooms() = rid %q not found`, rid)
	}

	r, _ := room.GetRoom(rid)
	if rooms[rid] != r {
		t.Fatalf(`TestGetRooms() = %v, want %v`, rooms[rid], r)
	}
}

func TestUserJoinRoom(t *testing.T) {
	uid := setupUser(t)
	rid := setupRoom(t)

	url := fmt.Sprintf("%s/users/%s/room", testMainURL, uid)

	_, err := sendRequest(url, http.MethodPut, []byte(rid))
	if err != nil {
		t.Fatalf(`TestUserJoinRoom() = %v, want nil`, err)
	}

	users, _ := lobby.GetUsersInRoom(rid)
	if users[0] != uid {
		t.Fatalf(`TestUserJoinRoom() = %v, want %v`, users[0], uid)
	}
}

func TestGetRoomUsers(t *testing.T) {
	uid, rid := setupUserInRoom(t)

	url := fmt.Sprintf("%s/rooms/%s", testMainURL, rid)

	resp, err := sendRequest(url, http.MethodGet, []byte{})
	if err != nil {
		t.Fatalf(`TestGetRoomUsers() = %v, want nil`, err)
	}

	var users []string
	err = json.Unmarshal([]byte(resp), &users)

	if len(users) != 1 || err != nil {
		t.Fatalf(`TestGetRoomUsers() = %v, %v, want [uuid], nil`, users, err)
	}

	if users[0] != uid {
		t.Fatalf(`TestGetRoomUsers()= %q, want %q`, users[0], uid)
	}

	real, _ := lobby.GetUsersInRoom(rid)
	if users[0] != real[0] {
		t.Fatalf(`TestGetUsersInRoom() = %v, want %v`, users[0], real[0])
	}
}

func TestGetGameTypes(t *testing.T) {
	resp, err := sendRequest(testMainURL+"/gameModes", http.MethodGet, []byte{})
	if err != nil {
		t.Fatalf(`TestGetGameTypes() = %v, want nil`, err)
	}

	var gTypes []string
	err = json.Unmarshal([]byte(resp), &gTypes)
	if len(gTypes) != 1 || err != nil {
		t.Fatalf(`TestGetGameTypes() = %v, %v, want [%s], nil`, gTypes, err, testType)
	}

	if gTypes[0] != testType {
		t.Fatalf(`TestGetGameTypes() %q, want %q`, gTypes[0], testType)
	}
}

func TestUpdateUserName(t *testing.T) {
	uid := setupUser(t)

	url := fmt.Sprintf("%s/users/%s/name", testMainURL, uid)

	_, err := sendRequest(url, http.MethodPut, []byte(altUser))
	if err != nil {
		t.Fatalf(`TestUpdateUserName() = %v, want nil`, err)
	}

	user, _ := user.GetUser(uid)
	if user.Name != altUser {
		t.Fatalf(`TestUpdateUserName() = %q, want %q`, user.Name, altUser)
	}
}

func TestUserLeaveRoom(t *testing.T) {
	uid, rid := setupUserInRoom(t)

	url := fmt.Sprintf("%s/users/%s/leave", testMainURL, uid)

	_, err := sendRequest(url, http.MethodPut, []byte(rid))
	if err != nil {
		t.Fatalf(`TestUserLeaveRoom() = %v, want nil`, err)
	}

	users, _ := lobby.GetUsersInRoom(rid)
	if len(users) > 0 {
		t.Fatalf(`TestUserLeaveRoom() = %v, want []`, users)
	}
}

func TestUpdateRoomName(t *testing.T) {
	rid := setupRoom(t)

	url := fmt.Sprintf("%s/rooms/%s/name", testMainURL, rid)

	_, err := sendRequest(url, http.MethodPut, []byte(altRoom))
	if err != nil {
		t.Fatalf(`TestUpdateRoomName() = %v, want nil`, err)
	}

	room, _ := room.GetRoom(rid)
	if room.Name != altRoom {
		t.Fatalf(`TestUpdateRoomName() = %q, want %q`, room.Name, altRoom)
	}
}

func TestUpdateRoomStatus(t *testing.T) {
	rid := setupRoom(t)

	url := fmt.Sprintf("%s/rooms/%s/status", testMainURL, rid)

	_, err := sendRequest(url, http.MethodPut, []byte("Updated"))
	if err != nil {
		t.Fatalf(`TestUpdateRoomStatus() = %v, want nil`, err)
	}

	room, _ := room.GetRoom(rid)
	if room.Status != "Updated" {
		t.Fatalf(`TestUpdateRoomStatus() = %q, want "Updated"`, room.Status)
	}
}

func TestUpdateRoomType(t *testing.T) {
	rid := setupRoom(t)

	url := fmt.Sprintf("%s/rooms/%s/type", testMainURL, rid)

	_, err := sendRequest(url, http.MethodPut, []byte(testType))
	if err != nil {
		t.Fatalf(`TestUpdateRoomType() = %v, want nil`, err)
	}

	room, _ := room.GetRoom(rid)
	if room.Type != testType {
		t.Fatalf(`TestUpdateRoomType() = %q, want %q`, room.Status, testType)
	}
}

func TestGetRoomURL(t *testing.T) {
	rid := setupReadyRoom(t)

	url := fmt.Sprintf("%s/rooms/%s/url", testMainURL, rid)

	resp, err := sendRequest(url, http.MethodGet, []byte{})
	if err != nil {
		t.Fatalf(`testGetRoomURL() = %v, want nil`, err)
	}

	if string(resp) != testConURL {
		t.Fatalf(`TestGetRoomURL() = %q, want %q`, resp, testConURL)
	}
}

func TestInitRoomGame(t *testing.T) {
	rid := setupReadyRoom(t)

	url := fmt.Sprintf("%s/rooms/%s/create", testMainURL, rid)

	_, err := sendRequest(url, http.MethodPut, []byte{})
	if err != nil {
		t.Fatalf(`TestInitRoomGame() = %v, want nil`, err)
	}
}

func TestUpdateRoomRules(t *testing.T) {
	rid := setupActiveRoom(t)

	url := fmt.Sprintf("%s/rooms/%s/rules", testMainURL, rid)

	_, err := sendRequest(url, http.MethodPut, []byte(newRules))
	if err != nil {
		t.Fatalf(`TestUpdateRoomRules() = %v, want nil`, err)
	}
}

func TestStartRoomGame(t *testing.T) {
	rid := setupActiveRoom(t)

	url := fmt.Sprintf("%s/rooms/%s/start", testMainURL, rid)

	_, err := sendRequest(url, http.MethodPut, []byte{})
	if err != nil {
		t.Fatalf(`TestStartRoomGame() = %v, want nil`, err)
	}
}

func TestEndRoomGame(t *testing.T) {
	rid := setupActiveRoom(t)

	url := fmt.Sprintf("%s/rooms/%s/start", testMainURL, rid)

	sendRequest(url, http.MethodPut, []byte{})

	url = fmt.Sprintf("%s/rooms/%s/end", testMainURL, rid)

	_, err := sendRequest(url, http.MethodPut, []byte{})
	if err != nil {
		t.Fatalf(`TestEndRoomGame() = %v, want nil`, err)
	}
}

func TestDeleteUser(t *testing.T) {
	uid := setupUser(t)

	url := fmt.Sprintf("%s/users/%s", testMainURL, uid)

	_, err := sendRequest(url, http.MethodDelete, []byte{})
	if err != nil {
		t.Fatalf(`TestDeleteUser() = %v, want nil`, err)
	}

	_, err = user.GetUser(uid)
	if err == nil {
		t.Fatalf(`TestDeleteUser() = %v, want err`, err)
	}
}

func TestDeleteRoom(t *testing.T) {
	rid := setupRoom(t)

	url := fmt.Sprintf("%s/rooms/%s", testMainURL, rid)

	_, err := sendRequest(url, http.MethodDelete, []byte{})
	if err != nil {
		t.Fatalf(`TestDeleteRoom() = %v, want nil`, err)
	}

	_, err = room.GetRoom(rid)
	if err == nil {
		t.Fatalf(`TestDeleteRoom() = %v, want err`, err)
	}
}

func sendRequest(url string, method string, body []byte) (string, error) {
	reqBody := bytes.NewReader(body)

	request, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return "", err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}

	resBody, err := ioutil.ReadAll(response.Body)
	return string(resBody), err

}

func setupClientSuite() {
	go server.Serve(testMainPort)

	time.Sleep(200 * time.Millisecond)

	registry.RegisterGameType("Test", testConURL)

}

func cleanUpAfterSuite() {

}

func setupRoom(t *testing.T) string {
	rid, _ := room.CreateRoom(testRoom)

	t.Cleanup(func() { cleanUpRoom(rid) })

	return rid
}

func cleanUpRoom(rid string) {
	room.DeleteRoom(rid)
}

func setupReadyRoom(t *testing.T) string {
	rid := setupRoom(t)

	room.UpdateRoomType(rid, testType)

	return rid
}

func setupActiveRoom(t *testing.T) string {
	rid := setupReadyRoom(t)

	gameclient.CreateGame(rid, testConURL)

	return rid
}

func setupUser(t *testing.T) string {
	uid, _ := user.CreateUser(testUser)

	t.Cleanup(func() { cleanUpUser(uid) })

	return uid
}

func cleanUpUser(uid string) {
	user.DeleteUser(uid)
}

func setupUserInRoom(t *testing.T) (string, string) {
	uid := setupUser(t)
	rid := setupRoom(t)

	err := lobby.JoinUserToRoom(uid, rid)
	if err != nil {
		t.Fatalf(`Setting up user: %v`, err)
	}

	t.Cleanup(func() { cleanUpLobby(rid, uid) })

	return uid, rid
}

func cleanUpLobby(uid string, rid string) {
	lobby.RemoveUserFromRoom(uid, rid)
}
