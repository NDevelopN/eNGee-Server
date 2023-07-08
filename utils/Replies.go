package utils

import "fmt"

func ReplyError(msg GameMsg, err error) (GameMsg, error) {
	reply := GameMsg{
		Type:    "Error",
		GID:     msg.GID,
		UID:     msg.UID,
		Content: "There was an issue with the " + msg.Type + " request",
	}

	return reply, fmt.Errorf("error handling %v request; %v", msg.Type, err)
}

func ReplyACK(msg GameMsg) GameMsg {
	return GameMsg{
		Type: "ACK",
		GID:  msg.GID,
		UID:  msg.UID,
	}
}
