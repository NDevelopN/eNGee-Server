package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func Extract[O *PlayerInfo | *GameRules | *GameOp | *RemovePlr](r *http.Request, obj O) error {
	reqBody, err := io.ReadAll(r.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(reqBody, obj)
	return err
}

func PackSend[O PlayerInfo | GameInfo | GameList | GameStatus | GameRules | Leader | string](w http.ResponseWriter, msg O) error {
	response, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	return err

}

func Ack(w http.ResponseWriter, msg string) error {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(msg))
	return err
}
