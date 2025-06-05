package document

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	docs := router.Group("/documents")
	{
		docs.POST("/", CreateDocumentController)
		docs.GET("/:id", GetDocumentByIdController)
		docs.GET("/user/:user_id", GetDocumentsByUserController)
		docs.PUT("/:id", UpdateDocumentController)
		docs.DELETE("/:id", DeleteDocumentController)
	}
}
