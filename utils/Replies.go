package utils

import (
	"encoding/json"
)

var ErrMsg = GameMsg{
	Type:    "Response",
	Content: `{cause: "Error", message: "unknown error"}`,
}

func CreateReply(msg GameMsg, cause string, message string) (GameMsg, error) {
	content, err := json.Marshal(Response{Cause: cause, Message: message})
	if err != nil {
		return ErrMsg, err
	}

	return GameMsg{
		UID:     msg.UID,
		GID:     msg.GID,
		Type:    "Response",
		Content: string(content),
	}, nil
}
