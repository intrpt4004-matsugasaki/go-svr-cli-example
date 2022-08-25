package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.GET("/time", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"time":    time.Now().Format("03:04:05"),
			"message": "hello client!",
		})
	})
	engine.Run(":3000")
}
