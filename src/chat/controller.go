package chat

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProcessPromptController(c *gin.Context) {
	var request PromptRequest

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

func ProcessPromptWithFileController(c *gin.Context) {
	var request FilePromptRequest

	request.Prompt = c.PostForm("prompt")
	if request.Prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'prompt' es requerido en el formulario."})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No se subió ningún archivo con el nombre de campo 'file'."})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al obtener el archivo: " + err.Error()})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al abrir el archivo subido: " + err.Error()})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer el contenido del archivo: " + err.Error()})
		return
	}

	mimeType := determineMIMETypeService(fileHeader.Filename)
	if mimeType == "" {
		mimeType = fileHeader.Header.Get("Content-Type")
		if mimeType == "" {
			log.Println("No se pudo determinar el MIME type del archivo, usando application/octet-stream")
			mimeType = "application/octet-stream" // Un genérico si no se sabe
		}
	}
	log.Printf("Archivo subido: %s, MIME Type detectado/usado: %s, Tamaño: %d bytes", fileHeader.Filename, mimeType, len(fileBytes))

	fileDataPart := &GeminiFileData{
		MIMEType: mimeType,
		Data:     fileBytes,
	}

	geminiTextResponse, err := ProcessPromptWithFileService(c.Request.Context(), request.Prompt, fileDataPart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GeminiResponse{Response: geminiTextResponse})
}
