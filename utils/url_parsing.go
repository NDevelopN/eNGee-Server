package utils

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func GetRequestIDs(request *http.Request) []string {
	splitPath := strings.Split(request.URL.Path, "/")

	ids := make([]string, 0)

	for _, pathPart := range splitPath {
		_, err := uuid.Parse(pathPart)
		if err == nil {
			ids = append(ids, pathPart)
		}
	}

	return ids
}
