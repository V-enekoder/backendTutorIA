package main

import (
	//"fmt"
	"github.com/V-enekoder/backendTutorIA/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	//"os"
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
	r.Run()
}
