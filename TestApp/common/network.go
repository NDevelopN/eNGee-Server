package common

import (
	"Engee-Server/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

const url = "http://localhost:8090/"

func PostUser(t *testing.T, user utils.User) (utils.User, error) {
	postBody, _ := json.Marshal(user)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(url+"users", "application/json", responseBody)
	if err != nil {
		t.Fatalf("PostUser failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("PostUser: Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return user, fmt.Errorf(string(body))
	}

	var nUser utils.User
	err = json.Unmarshal((body), &nUser)
	if err != nil {
		t.Fatalf("PostUser: Could not parse reply: %v", err)
	}

	return nUser, nil
}

func GetUser(t *testing.T, uid string) (utils.User, error) {
	resp, err := http.Get(url + "users/" + uid)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("GetUser: Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return utils.User{}, fmt.Errorf(string(body))
	}

	var nUser utils.User
	err = json.Unmarshal((body), &nUser)
	if err != nil {
		t.Fatalf("GetUser: Could not parse reply: %v", err)
	}

	return nUser, nil
}

func PutUser(t *testing.T, user utils.User) (utils.Response, error) {
	putBody, _ := json.Marshal(user)
	responseBody := bytes.NewBuffer(putBody)
	req, err := http.NewRequest(http.MethodPut, url+"users/"+user.UID, responseBody)

	if err != nil {
		t.Fatalf("PutUser: Request failed: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PutUser failed: %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("PutUser: Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return utils.Response{}, fmt.Errorf(string(body))
	}

	var reply utils.Response
	err = json.Unmarshal((body), &reply)
	if err != nil {
		t.Fatalf("PutUser: Could not parse reply: %v", err)
	}

	return reply, nil
}

func DeleteUser(t *testing.T, uid string) (utils.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url+"users/"+uid, bytes.NewBuffer([]byte(uid)))
	if err != nil {
		t.Fatalf("DeleteUser: Request failed: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("DeleteUser: Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return utils.Response{}, fmt.Errorf(string(body))
	}

	var reply utils.Response
	err = json.Unmarshal(body, &reply)
	if err != nil {
		t.Fatalf("DeletUser: Could not parse reply: %v", err)
	}

	return reply, nil

}

func PostGame(t *testing.T, game utils.Game) (utils.Game, error) {
	postBody, _ := json.Marshal(game)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(url+"games", "application/json", responseBody)
	if err != nil {
		t.Fatalf("PostGame failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("PostGame: Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return game, fmt.Errorf(string(body))
	}

	var nGame utils.Game
	err = json.Unmarshal((body), &nGame)
	if err != nil {
		t.Fatalf("PostGame: Could not parse reply: %v", err)
	}

	return nGame, nil
}

func GetGame(t *testing.T, gid string) (utils.Game, error) {
	resp, err := http.Get(url + "games/" + gid)
	if err != nil {
		t.Fatalf("GetGame failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("GetGame: Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return utils.Game{}, fmt.Errorf(string(body))
	}

	var nGame utils.Game
	err = json.Unmarshal((body), &nGame)
	if err != nil {
		t.Fatalf("GetGame: Could not parse reply: %v", err)
	}

	return nGame, nil

}

func PutGame(t *testing.T, game utils.Game) (utils.Response, error) {
	putBody, _ := json.Marshal(game)
	responseBody := bytes.NewBuffer(putBody)

	req, err := http.NewRequest(http.MethodPut, url+"games/"+game.GID, responseBody)
	if err != nil {
		t.Fatalf("PutGame: Request failed: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PutGame failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("PutGame: Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return utils.Response{}, fmt.Errorf(string(body))
	}

	var reply utils.Response
	err = json.Unmarshal(body, &reply)
	if err != nil {
		t.Fatalf("PutGame: Could not parse reply: %v", err)
	}

	return reply, nil
}

func DeleteGame(t *testing.T, gid string) (utils.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url+"games/"+gid, bytes.NewBuffer([]byte(gid)))
	if err != nil {
		t.Fatalf("DeleteGame: Request failed: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DeleteGame failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("DeleteGame: Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return utils.Response{}, fmt.Errorf(string(body))
	}

	reply := utils.Response{Cause: "", Message: ""}
	err = json.Unmarshal((body), &reply)
	if err != nil {
		t.Fatalf("DeleteGame: Could not parse reply: %v", err)
	}

	return reply, nil

}
