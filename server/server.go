package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	gameClient "Engee-Server/gameClient"
	registry "Engee-Server/gameRegistry"
	"Engee-Server/lobby"
	"Engee-Server/room"
	"Engee-Server/user"
	"Engee-Server/utils"
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

	router.POST("/users", postUser)
	router.POST("/rooms", postRoom)

	router.GET("/rooms", getRooms)
	router.GET("/rooms/:rid", getRoomUsers)
	router.GET("/rooms/:rid/url", getRoomURL)

	router.GET("/gameModes", getGameTypes)

	router.PUT("/users/:uid/name", updateUserName)
	router.PUT("/users/:uid/room", userJoinRoom)
	router.PUT("/users/:uid/leave", userLeaveRoom)

	router.PUT("/rooms/:rid/name", updateRoomName)
	router.PUT("/rooms/:rid/status", updateRoomStatus)
	router.PUT("/rooms/:rid/type", updateRoomType)
	router.PUT("/rooms/:rid/rules", updateRoomRules)

	router.PUT("/rooms/:rid/create", initRoomGame)
	router.PUT("/rooms/:rid/start", startRoomGame)
	router.PUT("/rooms/:rid/end", endRoomGame)

	router.DELETE("/users/:uid", deleteUser)
	router.DELETE("/rooms/:rid", deleteRoom)

	router.Run(":" + port)
}

func postUser(c *gin.Context) {
	reqBody, w := processMessage(c)
	uid, err := user.CreateUser(string(reqBody))

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Creating user: %v", err)
		return
	}

	sendReply(w, uid, http.StatusAccepted)
}

func postRoom(c *gin.Context) {
	reqBody, w := processMessage(c)
	rid, err := room.CreateRoom(string(reqBody))

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create room: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Creating room: %v", err)
		return
	}

	sendReply(w, rid, http.StatusAccepted)
}

func getRooms(c *gin.Context) {
	_, w := processMessage(c)
	rooms, err := room.GetRooms()

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get rooms: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Getting rooms: %v", err)
		return
	}

	roomsJSON, err := json.Marshal(rooms)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to package room info: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Marshalling rooms: %v", err)
		return
	}

	sendReply(w, string(roomsJSON), http.StatusOK)
}

func getRoomUsers(c *gin.Context) {
	_, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)

	users, err := lobby.GetUsersInRoom(ids[0])

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get room users: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Getting room users: %v", err)
		return
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to package room user info: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Marshalling room users: %v", err)
		return
	}

	sendReply(w, string(usersJSON), http.StatusOK)
}

func getRoomURL(c *gin.Context) {
	_, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)

	url, err := room.GetRoomURL(ids[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get room URL: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Getting room URL: %v", err)
		return
	}

	sendReply(w, url, http.StatusOK)
}

func getGameTypes(c *gin.Context) {
	_, w := processMessage(c)

	gTypes := registry.GetGameTypes()

	gTypesJSON, err := json.Marshal(gTypes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to package game types: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Marshalling game types: %v", err)
		return
	}

	sendReply(w, string(gTypesJSON), http.StatusOK)
}

func updateUserName(c *gin.Context) {
	reqBody, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)
	err := user.UpdateUserName(ids[0], string(reqBody))

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update user name: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Updating user name: %v", err)
		return
	}
}

func userJoinRoom(c *gin.Context) {
	reqBody, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)
	err := lobby.JoinUserToRoom(ids[0], string(reqBody))

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add uer to room: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Adding user to room: %v", err)
		return
	}
}

func userLeaveRoom(c *gin.Context) {
	reqBody, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)
	err := lobby.RemoveUserFromRoom(ids[0], string(reqBody))

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to remove user from room: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Removign user from room")
	}
}

func updateRoomName(c *gin.Context) {
	reqBody, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)
	err := room.UpdateRoomName(ids[0], string(reqBody))

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update room name: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Updating room name: %v", err)
		return
	}
}

func updateRoomStatus(c *gin.Context) {
	reqBody, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)
	err := room.UpdateRoomStatus(ids[0], string(reqBody))

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update room status: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Updating room status: %v", err)
		return
	}
}

func updateRoomType(c *gin.Context) {
	reqBody, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)
	err := room.UpdateRoomType(ids[0], string(reqBody))

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update room type: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Updating room type: %v", err)
		return
	}
}

func updateRoomRules(c *gin.Context) {
	reqBody, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)
	err := gameClient.SetGameRules(ids[0], string(reqBody))

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update room rules: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Updating room rules: %v", err)
		return
	}
}

func initRoomGame(c *gin.Context) {
	_, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)
	err := room.CreateRoomInstance(ids[0])

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update user name: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Updating user name: %v", err)
		return
	}
}

func startRoomGame(c *gin.Context) {
	_, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)
	err := gameClient.StartGame(ids[0])

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update user name: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Updating user name: %v", err)
		return
	}
}

func endRoomGame(c *gin.Context) {
	reqBody, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)
	err := user.UpdateUserName(ids[0], string(reqBody))

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update user name: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Updating user name: %v", err)
		return
	}
}

func deleteUser(c *gin.Context) {
	_, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)
	err := user.DeleteUser(ids[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete user: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Deleting user: %v", err)
		return
	}
}

func deleteRoom(c *gin.Context) {
	_, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)

	err := room.DeleteRoom(ids[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete room: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Deleting room: %v", err)
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
