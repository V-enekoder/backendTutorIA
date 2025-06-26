package chat

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	chat := router.Group("/chat")
	{
		chat.POST("/prompt/:id/:user_id", ProcessPromptController)
		chat.POST("/prompt/:id//:user_id/file", ProcessPromptWithFileController)
		chat.GET("/history/:id/:user_id", GetChatHistoryController)
	}
}
