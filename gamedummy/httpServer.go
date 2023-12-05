package gamedummy

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateMessage struct {
	RID    string `json:"rid"`
	Update string `json:"update"`
}

func CORSMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.Writer.WriteHeader(http.StatusOK)
			c.Writer.Write([]byte(""))
		}

		c.Next()
	}
}

func Serve(port string) {
	router := gin.Default()

	router.Use(CORSMiddleWare())

	router.POST("/games", postGame)
	router.PUT("/games/:id", updateGame)
	router.DELETE("/games/:id", deleteGame)

	router.Run(":" + port)
}

func postGame(c *gin.Context) {
	reqBody, w := processMessage(c)

	pURL, err := CreateNewInstance(string(reqBody))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create game: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Creating game: %v", err)
		return
	}

	err = sendReply(w, pURL, 200)
	if err != nil {
		log.Printf("[Error] Replying after creating game: %v", err)
		return
	}
}

func updateGame(c *gin.Context) {
	reqBody, w := processMessage(c)

	var um updateMessage

	err := json.Unmarshal(reqBody, &um)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse update message: %v", reqBody), http.StatusBadRequest)
		log.Printf("[Error] Reading update message: %v", err)
		return
	}

	switch um.Update {
	case ("Start"):
		err = StartInstance(um.RID)
	default:
		http.Error(w, fmt.Sprintf("Failed to update game, command: %q not recognised", um.Update), http.StatusBadRequest)
		log.Printf("[Error] Invalid update command %q", um.Update)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update game: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Updating(%s) game: %v", um.Update, err)
		return
	}

}

func deleteGame(c *gin.Context) {
	reqBody, w := processMessage(c)

	err := DeleteInstance(string(reqBody))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete game: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Deleting game: %v", err)
		return
	}

	err = sendReply(w, "", 200)
	if err != nil {
		log.Printf("[Error] Replying after deleting game: %v", err)
		return
	}
}

func processMessage(c *gin.Context) ([]byte, http.ResponseWriter) {
	w := c.Writer
	r := c.Request

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		log.Printf("[Error] Reading request body: %v", err)
		return nil, nil
	}

	return reqBody, w
}

func sendReply(w http.ResponseWriter, msg string, code int) error {
	w.WriteHeader(code)
	_, err := w.Write([]byte(msg))
	if err != nil {
		http.Error(w, "Could not write response", http.StatusInternalServerError)
		return fmt.Errorf("could not write response: %v", err)
	}

	return nil
}
