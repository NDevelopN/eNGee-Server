package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

func ReplyWarning(msg GameMsg, warning string) (GameMsg, error) {
	issue := Response{
		Cause:   "Warning",
		Message: warning,
	}

	var content string
	c, e := json.Marshal(issue)
	if e != nil {
		content = "Could not prepare warning message"
		e = fmt.Errorf("could not marshal warning message: %v", e)
	} else {
		content = string(c)
	}

	reply := GameMsg{
		Type:    "Response",
		GID:     msg.GID,
		UID:     msg.UID,
		Content: content,
	}

	return reply, e
}

func ReplyError(msg GameMsg, err error) (GameMsg, error) {
	issue := Response{
		Cause:   "Error",
		Message: err.Error(),
	}

	var content string
	c, e := json.Marshal(issue)
	if e != nil {
		content = "Could not prepare error message"
		e = fmt.Errorf("could not marshal error message: %v", e)
	} else {
		content = string(c)
	}

	reply := GameMsg{
		Type:    "Response",
		GID:     msg.GID,
		UID:     msg.UID,
		Content: content,
	}

	return reply, fmt.Errorf("error handling %v request: %v (%v)", msg.Type, err, e)
}

func ReplyACK(msg GameMsg, response string) (GameMsg, error) {
	accept := Response{
		Cause:   "Accept",
		Message: response,
	}

	var content string
	c, e := json.Marshal(accept)
	if e != nil {
		content = "Could not prepare accept message"
		log.Printf("[Error] Could not marshal accept message: %v", e)
	} else {
		content = string(c)
	}

	reply := GameMsg{
		Type:    "Response",
		GID:     msg.GID,
		UID:     msg.UID,
		Content: content,
	}

	return reply, e
}
