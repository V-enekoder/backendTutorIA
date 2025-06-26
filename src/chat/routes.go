package chat

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	chat := router.Group("/chat")
	{
		chat.POST("/prompt/:id", ProcessPromptController)
		chat.POST("/prompt/:id/file", ProcessPromptWithFileController)
	}
}
