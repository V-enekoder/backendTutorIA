package chat

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SendPromptController maneja la solicitud HTTP para el prompt de Gemini
func SendPromptController(c *gin.Context) {
	var request PromptRequest

	// BindJSON para validar y parsear el JSON de entrada
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de solicitud inválido: " + err.Error()})
		return
	}

	if request.Prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'prompt' no puede estar vacío"})
		return
	}

	// Llamar al servicio
	geminiTextResponse, err := ProcessPromptService(request.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Enviar la respuesta
	c.JSON(http.StatusOK, GeminiResponse{Response: geminiTextResponse})
}
