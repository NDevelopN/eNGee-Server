package gameclient

import (
	"Engee-Server/utils"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

var gameURLs = make(map[string]string)

func CreateGame(rid string, url string) error {
	if rid == "" {
		return fmt.Errorf("empty RID provided")
	}

	err := utils.ValidateURL(url)
	if err != nil {
		return err
	}

	_, found := gameURLs[rid]
	if found {
		return fmt.Errorf("game already exists for room %s", rid)
	}

	_, err = sendRequest(url+"/games", http.MethodPost, []byte(rid))

	if err != nil {
		return err
	}

	gameURLs[rid] = url + "/games/" + rid

	return nil
}

func EndGame(rid string) error {
	err := checkRID(rid)
	if err != nil {
		return err
	}

	url := gameURLs[rid]

	_, err = sendRequest(url, http.MethodDelete, []byte(rid))
	if err != nil {
		return err
	}

	delete(gameURLs, rid)
	return nil
}

func SetGameRules(rid string, rules string) error {
	err := checkRID(rid)
	if err != nil {
		return err
	}

	url := gameURLs[rid] + "/rules"

	_, err = sendRequest(url, http.MethodPut, []byte(rules))
	return err
}

func StartGame(rid string) error {
	err := checkRID(rid)
	if err != nil {
		return err
	}

	url := gameURLs[rid] + "/start"

	_, err = sendRequest(url, http.MethodPut, []byte{})
	return err
}

func PauseGame(rid string) error {
	err := checkRID(rid)
	if err != nil {
		return err
	}

	url := gameURLs[rid] + "/pause"

	_, err = sendRequest(url, http.MethodPut, []byte{})
	return err
}

func ResetGame(rid string) error {
	err := checkRID(rid)
	if err != nil {
		return err
	}

	url := gameURLs[rid] + "/reset"

	_, err = sendRequest(url, http.MethodPut, []byte{})
	return err
}

func RemovePlayer(rid string, targetUID string) error {
	err := checkRID(rid)
	if err != nil {
		return err
	}

	url := gameURLs[rid] + "/players/" + targetUID

	_, err = sendRequest(url, http.MethodDelete, []byte{})
	return err
}

func checkRID(rid string) error {
	_, found := gameURLs[rid]
	if !found {
		return fmt.Errorf("game %s not set up with URL", rid)
	}

	return nil
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
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to complete request: %s", resBody)

	}

	return string(resBody), nil
}
