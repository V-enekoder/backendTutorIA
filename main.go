package main

import (
	"net/http"

	"github.com/V-enekoder/backendTutorIA/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/V-enekoder/backendTutorIA/src/chat"
	"github.com/V-enekoder/backendTutorIA/src/document"
	"github.com/V-enekoder/backendTutorIA/src/project"
	"github.com/V-enekoder/backendTutorIA/src/user"
)

func init() {
	config.LoadEnv()
	config.ConnectDB()
	config.SyncDB()
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	user.RegisterRoutes(r)
	project.RegisterRoutes(r)
	chat.RegisterRoutes(r)
	document.RegisterRoutes(r)
	r.Run()
}
