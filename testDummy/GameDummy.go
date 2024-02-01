package testDummy

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func corsMiddleWare() gin.HandlerFunc {
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

	router.Use(corsMiddleWare())

	router.POST("/games", func(c *gin.Context) { c.Writer.Write([]byte{}) })
	router.PUT("/games/:id/start", func(c *gin.Context) { c.Writer.Write([]byte{}) })
	router.PUT("/games/:id/pause", func(c *gin.Context) { c.Writer.Write([]byte{}) })
	router.PUT("/games/:id/reset", func(c *gin.Context) { c.Writer.Write([]byte{}) })
	router.PUT("/games/:id/rules", func(c *gin.Context) { c.Writer.Write([]byte{}) })
	router.DELETE("/games/:id/players/:id", func(c *gin.Context) { c.Writer.Write([]byte{}) })
	router.DELETE("/games/:id", func(c *gin.Context) { c.Writer.Write([]byte{}) })

	router.Run(":" + port)
}

func sendReply(c *gin.Context, response []byte) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(response)
}
