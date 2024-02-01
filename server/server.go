package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

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

	router.POST("/users/:id", userHeartbeat)

	router.GET("/rooms", getRooms)
	router.GET("/rooms/:rid/users", getRoomUsers)
	router.GET("/rooms/:rid", getRoomInfo)

	router.GET("/gameModes", getGameModes)
	router.POST("/gameModes", postGameMode)
	router.POST("/gameModes/:gameMode", gameModeHeartbeat)

	router.PUT("/users/:uid/name", updateUserName)
	router.PUT("/users/:uid/room", userJoinRoom)
	router.PUT("/users/:uid/leave", userLeaveRoom)

	router.PUT("/rooms/:rid/name", updateRoomName)
	router.PUT("/rooms/:rid/status", updateRoomStatus)
	router.PUT("/rooms/:rid/mode", updateRoomGameMode)
	router.PUT("/rooms/:rid/rules", updateRoomRules)

	router.PUT("/rooms/:rid/create", initRoomGame)
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

	err = sendSimpleReply(w, "POST user", uid, http.StatusOK)
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
	}
}

func postRoom(c *gin.Context) {
	reqBody, w := processMessage(c)
	rid, err := room.CreateRoom(reqBody)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create room: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Creating room: %v", err)
		return
	}

	err = sendSimpleReply(w, "POST room", rid, http.StatusOK)
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
	}
}

func userHeartbeat(c *gin.Context) {
	_, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)

	err := user.Heartbeat(ids[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Hearbeat failed: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Receiving user heartbeat: %v", err)
		return
	}

	err = sendAccept(w, "HEARTBEAT user")
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
	}
}

func getRooms(c *gin.Context) {
	_, w := processMessage(c)
	rooms := room.GetRooms()

	roomsJSON, err := json.Marshal(rooms)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to package room info: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Marshalling rooms: %v", err)
		return
	}

	err = sendReply(w, "GET rooms", roomsJSON, http.StatusOK)
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
	}
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

	err = sendReply(w, "GET room/users", usersJSON, http.StatusOK)
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
	}
}

func getRoomInfo(c *gin.Context) {
	_, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)

	roomInfo, err := room.GetRoom(ids[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get room URL: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Getting room URL: %v", err)
		return
	}

	rInfo, err := json.Marshal(roomInfo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to package room info: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Marshaling room info: %v", err)
		return
	}

	err = sendReply(w, "GET room/info", rInfo, http.StatusOK)
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
		return
	}
}

func getGameModes(c *gin.Context) {
	_, w := processMessage(c)

	gameModes := registry.GetGameModes()

	gameModesJSON, err := json.Marshal(gameModes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to package game modes: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Marshalling game modes: %v", err)
		return
	}

	err = sendReply(w, "GET gamemodes", gameModesJSON, http.StatusOK)
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
	}
}

func postGameMode(c *gin.Context) {
	reqBody, w := processMessage(c)

	type stringPair struct {
		First  string `json:"first"`
		Second string `json:"second"`
	}

	var gameMode stringPair
	err := json.Unmarshal(reqBody, &gameMode)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal game mode: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Unmarshalling game mode: %v", err)
		return
	}

	err = registry.RegisterGameMode(gameMode.First, gameMode.Second)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update game mode: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Updating game mode: %v", err)
		return
	}

	err = sendAccept(w, "POST gamemode")
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
	}
}

func gameModeHeartbeat(c *gin.Context) {
	_, w := processMessage(c)
	splitPath := strings.Split(c.Request.URL.Path, "/")
	modeName := splitPath[len(splitPath)-1]

	err := registry.Heartbeat(modeName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to accept heartbeat: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Receiving gamemode heartbeat: %v", err)
		return
	}

	err = sendAccept(w, "HEARTBEAT gamemode")
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
	}
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

	err = sendAccept(w, "PUT user/name")
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
	}
}

func userJoinRoom(c *gin.Context) {
	reqBody, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)
	err := lobby.JoinUserToRoom(ids[0], string(reqBody))

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add user to room: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Adding user to room: %v", err)
		return
	}

	err = sendAccept(w, "PUT user/room")
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
	}
}

func userLeaveRoom(c *gin.Context) {
	reqBody, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)
	err := lobby.RemoveUserFromRoom(ids[0], string(reqBody))

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to remove user from room: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Removing user from room: %v", err)
		return
	}

	err = sendAccept(w, "PUT user/leave")
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
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

	err = sendAccept(w, "PUT room/name")
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
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

	err = sendAccept(w, "PUT room/status")
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
		return
	}
}

func updateRoomGameMode(c *gin.Context) {
	reqBody, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)
	err := room.UpdateRoomGameMode(ids[0], string(reqBody))

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update room game mode: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Updating room game mode: %v", err)
		return
	}

	err = sendAccept(w, "PUT room/gamemode")
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
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

	err = sendAccept(w, "PUT room/rules")
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
		return
	}
}

func initRoomGame(c *gin.Context) {
	_, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)
	err := room.InitializeRoomGame(ids[0])

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update user name: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Updating user name: %v", err)
		return
	}

	err = sendAccept(w, "PUT room/init")
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
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

	err = sendAccept(w, "PUT room/end")
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
		return
	}
}

func deleteUser(c *gin.Context) {
	_, w := processMessage(c)
	ids := utils.GetRequestIDs(c.Request)

	err := lobby.RemoveUserFromAllRooms(ids[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete room(s) user is in: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Removing deleting user from room(s): %v", err)
		//No return, want to complete deleting user regardless
	}

	err = user.DeleteUser(ids[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete user: %v", err), http.StatusInternalServerError)
		log.Printf("[Error] Deleting user: %v", err)
		return
	}

	err = sendAccept(w, "DELETE user")
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
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

	err = sendAccept(w, "DELETE room")
	if err != nil {
		log.Printf("[Error] Sending reply: %v", err)
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

func sendReply(w http.ResponseWriter, request string, msg []byte, code int) error {
	w.WriteHeader(code)
	_, err := w.Write(msg)
	if err != nil {
		http.Error(w, "Could not write response", http.StatusInternalServerError)
		return fmt.Errorf("could not write %s response: %v", request, err)
	}

	return nil
}

func sendSimpleReply(w http.ResponseWriter, request string, msg string, code int) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not marshal message: %w", err)
	}

	return sendReply(w, request, msgJSON, code)
}

func sendAccept(w http.ResponseWriter, request string) error {
	return sendReply(w, request, []byte{}, http.StatusAccepted)
}
