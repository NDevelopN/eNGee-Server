package user

import (
	"time"

	"github.com/google/uuid"

	sErr "Engee-Server/stockErrors"
	"Engee-Server/utils"
)

type User struct {
	UID    string `json:"uid"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

var users = make(map[string]User)
var heartbeats map[string]time.Time

func CreateUser(name string) (string, error) {
	if name == "" {
		return "", &sErr.EmptyValueError{
			Field: "Name",
		}
	}

	var newUser User
	newUser.UID = uuid.NewString()
	newUser.Name = name
	newUser.Status = "New"

	if heartbeats == nil {
		heartbeats = make(map[string]time.Time)
		go utils.MonitorHeartbeats(&heartbeats, DeleteUser)
	}

	users[newUser.UID] = newUser
	heartbeats[newUser.UID] = time.Now()

	return newUser.UID, nil
}

func Heartbeat(uid string) error {
	_, err := GetUser(uid)
	if err != nil {
		return err
	}

	heartbeats[uid] = time.Now()

	return nil
}

func GetUser(uid string) (User, error) {
	user, found := users[uid]
	if !found {
		return user, &sErr.MatchNotFoundError[string]{
			Space: "Users",
			Field: "UID",
			Value: uid,
		}
	}

	return user, nil
}

func UpdateUserName(uid string, name string) error {
	if name == "" {
		return &sErr.EmptyValueError{
			Field: "Name",
		}
	}

	user, err := GetUser(uid)
	if err != nil {
		return err
	}

	user.Name = name
	users[uid] = user

	return nil
}

func UpdateUserStatus(uid string, status string) error {
	if status == "" {
		return &sErr.EmptyValueError{
			Field: "Status",
		}
	}

	user, err := GetUser(uid)
	if err != nil {
		return err
	}

	user.Status = status
	users[uid] = user

	return nil
}

func DeleteUser(uid string) error {
	_, err := GetUser(uid)
	if err != nil {
		return err
	}

	delete(users, uid)
	delete(heartbeats, uid)

	return nil
}
