package project

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	projects := router.Group("/projects")
	{
		projects.POST("/", CreateProjectController)          

		projects.GET("/:id", GetProjectByIdController)      

		projects.PUT("/:id", UpdateProjectController)    
		  
		projects.DELETE("/:id", DeleteProjectByIdController) 
	}
}