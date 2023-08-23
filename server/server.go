package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	db "Engee-Server/database"
	g "Engee-Server/game"
	gamespace "Engee-Server/gamespace"
	u "Engee-Server/user"
	"Engee-Server/utils"
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

	router.GET("/users/:id", getUsers)
	router.POST("/users", postUsers)
	router.PUT("/users/:id", putUsers)
	router.DELETE("/users/:id", deleteUsers)

	router.GET("/types", getTypes)

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

func reply[m utils.Message](w http.ResponseWriter, msg m, code int) error {
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

	err = reply(w, utils.Response{Cause: "Accept", Message: "Game deleted successfully"}, http.StatusAccepted)
	if err != nil {
		log.Printf("[Error] Replying: %v", err)
	}
}

func postGames(c *gin.Context) {
	reqBody, w := intake(c)

	var game utils.Game
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

	//TODO
	game.GID = gid
	game.CurPlrs = 1

	msg := utils.GameMsg{
		Type: "Init",
		UID:  game.Leader,
		GID:  game.GID,
	}

	user, err := u.GetUser(game.Leader)
	if err != nil {
		log.Printf("[Error] Failed to get user matching leader id: %v", err)
		err = g.DeleteGame(gid)
		if err != nil {
			log.Printf("[Error] Failed to delete game with no leader: %v", err)
		}

		http.Error(w, "Failed to get leader", http.StatusBadRequest)
		return
	}

	user.GID = game.GID
	err = u.UpdateUser(user)
	if err != nil {
		log.Printf("[Error] Failed to get update leader user: %v", err)
		err = g.DeleteGame(gid)
		if err != nil {
			log.Printf("[Error] Failed to delete game with invalid leader: %v", err)
		}
		http.Error(w, "Failed to update leader", http.StatusInternalServerError)
		return
	}

	go func() {
		_, err = gamespace.GamespaceHandle(msg)
		if err != nil {
			log.Printf("[Error] Failed to initialize game: %v", err)
			err = g.DeleteGame(gid)
			if err != nil {
				log.Printf("[Error] Failed to delete game with failed Init: %v", err)
			}
		}
	}()

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

	var game utils.Game
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

	err = reply(w, utils.Response{Cause: "Accept", Message: "Game updated successfully"}, http.StatusAccepted)
	if err != nil {
		log.Printf("[Error] Replying: %v", err)
	}
}

func getUsers(c *gin.Context) {
	w := c.Writer

	uid, err := GetID(c)
	if err != nil {
		http.Error(w, "Could not get user ID from request path", http.StatusBadRequest)
		log.Printf("[Error] Getting user ID: %v", err)
		return
	}

	user, err := u.GetUser(uid)
	if err != nil {
		http.Error(w, "Failed to get matching user", http.StatusInternalServerError)
		log.Printf("[Error] Getting user: %v", err)
		return
	}

	err = reply(w, user, http.StatusOK)
	if err != nil {
		log.Printf("[Error] Replying to GetUser: %v", err)
	}
}

func postUsers(c *gin.Context) {
	reqBody, w := intake(c)

	var user utils.User
	err := json.Unmarshal(reqBody, &user)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		log.Printf("[Error] Parsing request body: %v", err)
		return
	}

	uid, err := u.CreateUser(user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		log.Printf("[Error] Creating game: %v", err)
		return
	}

	//TODO
	user.UID = uid
	user.Status = "New"

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

	var user utils.User
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

	err = u.UpdateUser(user)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		log.Printf("[Error] Updating user: %v", err)
		return
	}

	err = reply(w, utils.Response{Cause: "Accept", Message: "User updated successfully"}, http.StatusAccepted)
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

	err = u.DeleteUser(uid)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		log.Printf("[Error] Deleting user: %v", err)
		return
	}

	err = reply(w, utils.Response{Cause: "Accept", Message: "User deleted successfully"}, http.StatusAccepted)
	if err != nil {
		log.Printf("[Error] Replying: %v", err)
	}
}

func getTypes(c *gin.Context) {
	w := c.Writer

	types, err := db.GetGameTypes()
	if err != nil {
		http.Error(w, "Failed to get types", http.StatusInternalServerError)
		log.Printf("[Error] Getting types: %v", err)
		return
	}

	err = reply(w, types, http.StatusOK)
	if err != nil {
		log.Printf("[Error] Replying to GetTypes: %v", err)
	}
}
