package gamedummy

import (
	"Engee-Server/utils"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	router.PUT("/games/:id/start", startGame)
	router.PUT("/games/:id/pause", pauseGame)
	router.PUT("/games/:id/reset", resetGame)
	router.PUT("/games/:id/rules", updateGameRules)
	router.DELETE("/games/:id/players/:id", removePlayer)
	router.DELETE("/games/:id", deleteGame)

	router.Run(":" + port)
}

func postGame(c *gin.Context) {
	reqBody, w := processMessage(c)

	err := CreateNewInstance(string(reqBody))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create game: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Creating game: %v", err)
		return
	}

	err = sendReply(w, "", 200)
	if err != nil {
		log.Printf("[Error] Replying after creating game: %v", err)
		return
	}
}

func startGame(c *gin.Context) {
	_, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)

	err := StartInstance(ids[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to start game: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Starting game: %v", err)
		return
	}

	err = sendReply(w, "", 200)
	if err != nil {
		log.Printf("[Error] Replying after starting game: %v", err)
		return
	}
}

func pauseGame(c *gin.Context) {
	_, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)

	err := PauseInstance(ids[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to pause/unpause game: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Pausing/Unpausing game: %v", err)
		return
	}

	err = sendReply(w, "", 200)
	if err != nil {
		log.Printf("[Error] Replying after pausing game: %v", err)
		return
	}
}

func resetGame(c *gin.Context) {
	_, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)

	err := ResetInstance(ids[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to reset game: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Resetting game: %v", err)
		return
	}

	err = sendReply(w, "", 200)
	if err != nil {
		log.Printf("[Error] Replying after resetting game: %v", err)
		return
	}
}

func updateGameRules(c *gin.Context) {
	reqBody, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)

	err := SetInstanceRules(ids[0], string(reqBody))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update game rules: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Updating game rules: %v", err)
		return
	}

	err = sendReply(w, "", 200)
	if err != nil {
		log.Printf("[Error] Replying after updating game rules: %v", err)
		return
	}
}

func removePlayer(c *gin.Context) {
	_, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)

	err := RemovePlayerFromInstance(ids[0], ids[1])
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to remove player from game: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Removing player from game: %v", err)
		return
	}

	err = sendReply(w, "", 200)
	if err != nil {
		log.Printf("[Error] Replying after removing player: %v", err)
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
