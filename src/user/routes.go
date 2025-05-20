package user

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	users := router.Group("/users")
	{
		users.POST("/", CreateUserController)

		users.GET("/id/:id", GetUserByIdController)

		users.PUT("/:id", UpdateUserController)                  // General user update
		users.PUT("/password/:id", UpdatePasswordUserController) // Password update

		users.DELETE("/:id", DeleteUserbyIdController)
	}
}
