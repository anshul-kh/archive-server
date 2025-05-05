package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var GlobalJobMapper *ExecutionMapper
var GlobalLangHandler *LangHandler
var GlobalCommander *Commander

func main() {
	GlobalJobMapper = NewExecutionMapper()
	GlobalLangHandler = NewLangHandler()
	GlobalCommander = NewCommander()

	commander := NewCommander()
	commander.GetListOfContainers()
	commander.InitServer()

	r := gin.Default()

	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all, or specify: []string{"http://localhost:3000"}
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "PONG",
		})
	})

	r.POST("/exec", executionHandler)
	r.GET("/job_result/:job_id", resultHandler)

	r.Run()
}
