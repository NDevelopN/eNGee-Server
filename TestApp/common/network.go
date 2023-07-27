package common

import (
	"Engee-Server/utils"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

const url = "http://localhost:8090/"

func PostUser(t *testing.T, user utils.User) (utils.User, utils.Issue) {
	postBody, _ := json.Marshal(user)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(url+"users", "application/json", responseBody)
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var nUser utils.User
	eMsg := utils.Issue{Cause: "", Message: ""}
	err = json.Unmarshal((body), &nUser)
	if err != nil {
		err = json.Unmarshal((body), &eMsg)
		if err != nil {
			t.Fatalf("Could not parse reply: %v", err)
		}
	}

	return nUser, eMsg
}

func GetUser(t *testing.T, uid string) (utils.User, utils.Issue) {
	resp, err := http.Get(url + "users/" + uid)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var nUser utils.User
	eMsg := utils.Issue{Cause: "", Message: ""}
	err = json.Unmarshal((body), &nUser)
	if err != nil {
		err = json.Unmarshal((body), &eMsg)
		if err != nil {
			t.Fatalf("Could not parse reply: %v", err)
		}
	}

	return nUser, eMsg

}

func PutUser(t *testing.T, user utils.User) (utils.User, utils.Issue) {
	putBody, _ := json.Marshal(user)
	responseBody := bytes.NewBuffer(putBody)
	resp, err := http.NewRequest(http.MethodPut, url+"users", responseBody)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var nUser utils.User
	eMsg := utils.Issue{Cause: "", Message: ""}
	err = json.Unmarshal((body), &nUser)
	if err != nil {
		err = json.Unmarshal((body), &eMsg)
		if err != nil {
			t.Fatalf("Could not parse reply: %v", err)
		}
	}

	return nUser, eMsg
}

func DeleteUser(t *testing.T, uid string) utils.Issue {
	resp, err := http.NewRequest(http.MethodDelete, url+"users/"+uid, bytes.NewBuffer([]byte(uid)))
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	eMsg := utils.Issue{Cause: "", Message: ""}
	err = json.Unmarshal((body), &eMsg)
	if err != nil {
		t.Fatalf("Could not parse reply: %v", err)
	}

	return eMsg

}

func PostGame(t *testing.T, game utils.Game) (utils.Game, utils.Issue) {
	postBody, _ := json.Marshal(game)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(url+"games", "application/json", responseBody)
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var nGame utils.Game
	eMsg := utils.Issue{Cause: "", Message: ""}
	err = json.Unmarshal((body), &nGame)
	if err != nil {
		err = json.Unmarshal((body), &eMsg)
		if err != nil {
			t.Fatalf("Could not parse reply: %v", err)
		}
	}

	return nGame, eMsg
}

func GetGame(t *testing.T, gid string) (utils.Game, utils.Issue) {
	resp, err := http.Get(url + "games/" + gid)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var nGame utils.Game
	eMsg := utils.Issue{Cause: "", Message: ""}
	err = json.Unmarshal((body), &nGame)
	if err != nil {
		err = json.Unmarshal((body), &eMsg)
		if err != nil {
			t.Fatalf("Could not parse reply: %v", err)
		}
	}

	return nGame, eMsg

}

func PutGame(t *testing.T, game utils.Game) (utils.Game, utils.Issue) {
	putBody, _ := json.Marshal(game)
	responseBody := bytes.NewBuffer(putBody)
	resp, err := http.NewRequest(http.MethodPut, url+"games", responseBody)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var nGame utils.Game
	eMsg := utils.Issue{Cause: "", Message: ""}
	err = json.Unmarshal((body), &nGame)
	if err != nil {
		err = json.Unmarshal((body), &eMsg)
		if err != nil {
			t.Fatalf("Could not parse reply: %v", err)
		}
	}

	return nGame, eMsg
}

func DeleteGame(t *testing.T, gid string) utils.Issue {
	resp, err := http.NewRequest(http.MethodDelete, url+"games/"+gid, bytes.NewBuffer([]byte(gid)))
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	eMsg := utils.Issue{Cause: "", Message: ""}
	err = json.Unmarshal((body), &eMsg)
	if err != nil {
		t.Fatalf("Could not parse reply: %v", err)
	}

	return eMsg

}
