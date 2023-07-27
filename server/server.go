package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	g "Engee-Server/game"
	gamespace "Engee-Server/gamespace"
	p "Engee-Server/user"
	u "Engee-Server/utils"
	utils "Engee-Server/utils"
)

func CORSMiddleware() gin.HandlerFunc {
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

func Serve() {
	router := gin.Default()

	router.Use(CORSMiddleware())

	router.GET("/games", getGames)
	router.POST("/games", postGames)
	router.PUT("/games/:id", putGames)
	router.DELETE("/games/:id", deleteGames)

	router.POST("/users", postUsers)
	router.PUT("/users/:id", putUsers)
	router.DELETE("/users/:id", deleteUsers)

	//This special case creates a websocket connection
	router.GET("/games/:id", Connect)

	router.Run("localhost:8090")
}

func intake(c *gin.Context) ([]byte, http.ResponseWriter) {
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

func GetID(c *gin.Context) (string, error) {
	r := c.Request

	path := r.URL.Path
	s := strings.Split(path, "/")

	id := s[len(s)-1]
	_, err := uuid.Parse(id)
	if err != nil {
		return "", fmt.Errorf("could not parse id: %v", err)
	}

	return id, nil
}

func reply[m u.Message](w http.ResponseWriter, msg m, code int) error {
	response, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, "Could not marshal response", http.StatusInternalServerError)
		return fmt.Errorf("could not marshal message: %v", err)
	}

	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, "Could not write response", http.StatusInternalServerError)
		return fmt.Errorf("could not write resposne: %v", err)
	}

	return nil
}

func getGames(c *gin.Context) {
	w := c.Writer

	games, err := g.GetGames()
	if err != nil {
		http.Error(w, "Could not get game ID from request path", http.StatusInternalServerError)
		log.Printf("[Error] Getting game list")
	}

	err = reply(w, games, http.StatusOK)
	if err != nil {
		log.Printf("[Error] Replying: %v", err)
	}
}

func deleteGames(c *gin.Context) {
	w := c.Writer

	gid, err := GetID(c)
	if err != nil {
		http.Error(w, "Could not get game ID from request path", http.StatusBadRequest)
		log.Printf("[Error] Getting game ID: %v", err)
	}

	err = g.DeleteGame(gid)
	if err != nil {
		http.Error(w, "Failed to delete game", http.StatusInternalServerError)
		log.Printf("[Error] Deleting game: %v", err)
		return
	}
	if err != nil {
		log.Printf("[Error] Replying: %v", err)
	}
}

func postGames(c *gin.Context) {
	reqBody, w := intake(c)

	var game u.Game
	err := json.Unmarshal(reqBody, &game)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		log.Printf("[Error] Parsing request body: %v", err)
		return
	}

	gid, err := g.CreateGame(game)
	if err != nil {
		http.Error(w, "Failed to create game", http.StatusInternalServerError)
		log.Printf("[Error] Creating game: %v", err)
		return
	}

	aborted := false

	go func() {
		time.Sleep(5 * time.Second)
		plrs, err := g.GetGamePlayers(gid)
		if len(plrs) == 0 || err != nil {
			log.Printf("[Error] no players found for game after creation: %v", err)
			g.DeleteGame(gid)
			aborted = true
		}
	}()

	if aborted {
		return
	}

	game.GID = gid

	msg := utils.GameMsg{
		Type: "Init",
		UID:  game.Leader,
		GID:  game.GID,
	}

	if aborted {
		return
	}

	user, err := p.GetUser(game.Leader)
	if err != nil {
		g.DeleteGame(gid)
		http.Error(w, "Failed to get leader", http.StatusBadRequest)
		log.Printf("[Error] Failed to get user matching leader id: %v", err)
		return
	}

	if aborted {
		return
	}

	user.GID = game.GID
	err = p.UpdateUser(user)
	if err != nil {
		http.Error(w, "Failed to update leader", http.StatusInternalServerError)
		log.Printf("[Error] Failed to get update leader user: %v", err)
		return
	}

	if aborted {
		return
	}

	_, err = gamespace.GamespaceHandle(msg)
	if err != nil {
		http.Error(w, "Failed to initialize gamespace", http.StatusInternalServerError)
		log.Printf("[Error] Failed to initialize game: %v", err)
		return
	}

	if aborted {
		return
	}
	err = reply(w, game, http.StatusCreated)
	if err != nil {
		log.Printf("[Error] Replying: %v", err)
	}

}

func putGames(c *gin.Context) {
	reqBody, w := intake(c)
	if reqBody == nil || w == nil {
		return
	}

	gid, err := GetID(c)
	if err != nil {
		http.Error(w, "Could not get game ID from request path", http.StatusBadRequest)
		log.Printf("[Error] Getting game ID: %v", err)
	}

	var game u.Game
	err = json.Unmarshal(reqBody, &game)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		log.Printf("[Error] Parsing request body: %v", err)
		return
	}

	if game.GID != gid {
		http.Error(w, "Mismatching request target and game ID", http.StatusBadRequest)
		log.Printf("[Error] Mismatching request target: %v, %v", gid, game.GID)
		return
	}

	err = g.UpdateGame(game)
	if err != nil {
		http.Error(w, "Failed to update game", http.StatusInternalServerError)
		log.Printf("[Error] Updating game: %v", err)
		return
	}

	err = reply(w, u.Response{Cause: "Accept", Message: "Game updated successfully"}, http.StatusAccepted)
	if err != nil {
		log.Printf("[Error] Replying: %v", err)
	}
}

func postUsers(c *gin.Context) {
	reqBody, w := intake(c)

	var user u.User
	err := json.Unmarshal(reqBody, &user)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		log.Printf("[Error] Parsing request body: %v", err)
		return
	}

	uid, err := p.CreateUser(user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		log.Printf("[Error] Creating game: %v", err)
		return
	}

	user.UID = uid

	err = reply(w, user, http.StatusCreated)
	if err != nil {
		log.Printf("[Error] Replying: %v", err)
	}
}

func putUsers(c *gin.Context) {

	reqBody, w := intake(c)
	if reqBody == nil || w == nil {
		return
	}

	uid, err := GetID(c)

	if err != nil {
		http.Error(w, "Could not get user ID from request path", http.StatusBadRequest)
		log.Printf("[Error] Getting user ID: %v", err)
	}

	var user u.User
	err = json.Unmarshal(reqBody, &user)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		log.Printf("[Error] Parsing request body: %v", err)
		return
	}

	if user.UID != uid {
		http.Error(w, "Mismatching request target and user ID", http.StatusBadRequest)
		log.Printf("[Error] Mismatching request target: %v, %v", uid, user.GID)
		return
	}

	err = p.UpdateUser(user)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		log.Printf("[Error] Updating user: %v", err)
		return
	}

	err = reply(w, u.Response{Cause: "Accept", Message: "User updated successfully"}, http.StatusAccepted)
	if err != nil {
		log.Printf("[Error] Replying: %v", err)
	}
}

func deleteUsers(c *gin.Context) {
	w := c.Writer

	uid, err := GetID(c)
	if err != nil {
		http.Error(w, "Could not get user ID from request path", http.StatusBadRequest)
		log.Printf("[Error] Getting user ID: %v", err)
	}

	err = p.DeleteUser(uid)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		log.Printf("[Error] Deleting user: %v", err)
		return
	}

	err = reply(w, u.Response{Cause: "Accept", Message: "User deleted successfully"}, http.StatusAccepted)
	if err != nil {
		log.Printf("[Error] Replying: %v", err)
	}
}
