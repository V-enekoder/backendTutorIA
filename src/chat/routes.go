package chat

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	chat := router.Group("/chat")
	{
		chat.POST("/prompt", ProcessPromptController)
		chat.POST("/prompt/file", ProcessPromptWithFileController)
	}
}
