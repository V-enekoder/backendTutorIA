package chat

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProcessPromptController(c *gin.Context) {
	var uriParams UriParams
	var promptBody PromptBody

	// PASO 1: Vincular y validar SOLO la URI.
	if err := c.ShouldBindUri(&uriParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de contexto inválido en la URL. Debe ser un número entre 1 y 20."})
		return
	}

	// PASO 2: Vincular y validar SOLO el cuerpo JSON.
	if err := c.ShouldBindJSON(&promptBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de la solicitud inválido o campo 'prompt' faltante."})
		return
	}

	// PASO 3: Llamar al servicio con los datos de ambos structs.
	geminiTextResponse, err := ProcessPromptService(promptBody.Prompt, uriParams.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GeminiResponse{Response: geminiTextResponse})
}

func ProcessPromptWithFileController(c *gin.Context) {
	var uriParams UriParams // Usamos el mismo struct para la URI

	// PASO 1: Vincular y validar SOLO la URI.
	if err := c.ShouldBindUri(&uriParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de contexto inválido en la URL. Debe ser un número entre 1 y 20."})
		return
	}

	// El resto del código para manejar el formulario se mantiene igual...
	prompt := c.PostForm("prompt")
	if prompt == "" {
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
			mimeType = "application/octet-stream"
		}
	}
	log.Printf("Archivo subido: %s, MIME Type: %s, Tamaño: %d bytes, ID de Contexto: %d", fileHeader.Filename, mimeType, len(fileBytes), uriParams.ID)

	fileDataPart := &GeminiFileData{
		MIMEType: mimeType,
		Data:     fileBytes,
	}

	geminiTextResponse, err := ProcessPromptWithFileService(c.Request.Context(), prompt, fileDataPart, uriParams.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GeminiResponse{Response: geminiTextResponse})
}
