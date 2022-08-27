package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthReq struct {
	Id string `form:"id"`
	Pw string `form:"pw"`
}

type ReverseReq struct {
	Text string `form:"text"`
}

const Token = "44gh348th34g93hg803g0"

func main() {
	engine := gin.Default()

	engine.POST("/auth", func(c *gin.Context) {
		req := AuthReq{}
		c.Bind(&req)

		if req.Id == "root" && req.Pw == "password" {
			c.JSON(http.StatusOK, gin.H{
				"token": Token,
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "failed.",
			})
		}

	})

	engine.GET("/time", func(c *gin.Context) {
		if c.Query("token") != Token {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "illegal token",
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"time": time.Now().Format("03:04:05"),
		})
	})

	ua := ""
	engine.Use(func(c *gin.Context) {
		ua = c.GetHeader("User-Agent")
		c.Next()
	})
	engine.GET("/user-agent", func(c *gin.Context) {
		if c.Query("token") != Token {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "illegal token",
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user-agent": ua,
		})
	})

	engine.POST("/reverse", func(c *gin.Context) {
		if c.Query("token") != Token {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "illegal token",
			})

			return
		}

		req := ReverseReq{}
		c.Bind(&req)

		var result []rune
		for i := len(req.Text) - 1; i >= 0; i-- {
			result = append(result, []rune(req.Text)[i])
		}

		c.JSON(http.StatusOK, gin.H{
			"text": string(result),
		})
	})

	engine.Run(":3000")
}
