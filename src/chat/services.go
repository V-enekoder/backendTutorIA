package chat

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/V-enekoder/backendTutorIA/config"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func ProcessPromptService(promptText string) (string, error) {
	apiKey := config.GetGeminiAPIKey()

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("Error creando el cliente de Gemini: %v", err)
		return "", errors.New("no se pudo inicializar el cliente de Gemini")
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	resp, err := model.GenerateContent(ctx, genai.Text(promptText))
	if err != nil {
		log.Printf("Error generando contenido con Gemini: %v", err)
		return "", errors.New("error al comunicarse con la API de Gemini")
	}

	var responseBuilder strings.Builder
	if resp != nil && len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				responseBuilder.WriteString(string(txt))
			}
		}
	} else {
		log.Println("Respuesta de Gemini vacía o inesperada")
		return "", errors.New("respuesta vacía o inesperada de Gemini")
	}

	return responseBuilder.String(), nil
}

func determineMIMETypeService(filename string) string {
	if strings.HasSuffix(filename, ".pdf") {
		return "application/pdf"
	} else if strings.HasSuffix(filename, ".png") {
		return "image/png"
	} else if strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") {
		return "image/jpeg"
	} else if strings.HasSuffix(filename, ".txt") {
		return "text/plain"
	}

	return ""
}

func ProcessPromptWithFileService(ctx context.Context, promptText string, fileData *GeminiFileData) (string, error) {
	if fileData == nil || len(fileData.Data) == 0 {
		return "", errors.New("ProcessPromptWithFileService requiere datos de archivo válidos")
	}
	if fileData.MIMEType == "" {
		return "", errors.New("ProcessPromptWithFileService requiere un MIMEType para el archivo")
	}

	apiKey := config.GetGeminiAPIKey()

	if ctx == nil {
		log.Println("Advertencia: ProcessPromptWithFileService recibió un contexto nulo, usando context.Background().")
		ctx = context.Background()
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("Error creando el cliente de Gemini: %v", err)
		return "", errors.New("no se pudo inicializar el cliente de Gemini")
	}
	defer client.Close()

	// "gemini-1.5-flash" o "gemini-1.5-pro" son recomendados para multimodalidad.
	model := client.GenerativeModel("gemini-1.5-flash")

	// Construir las partes del contenido
	var parts []genai.Part
	// Añadir la parte del archivo primero
	parts = append(parts, genai.Blob{MIMEType: fileData.MIMEType, Data: fileData.Data})
	// Luego el prompt de texto
	parts = append(parts, genai.Text(promptText))

	resp, err := model.GenerateContent(ctx, parts...)
	if err != nil {
		log.Printf("Error generando contenido con Gemini: %v", err)
		return "", errors.New("error al comunicarse con la API de Gemini")
	}

	var responseBuilder strings.Builder
	if resp != nil && len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				responseBuilder.WriteString(string(txt))
			}
		}
	} else {
		log.Println("Respuesta de Gemini vacía o inesperada")
		return "", errors.New("respuesta vacía o inesperada de Gemini")
	}

	return responseBuilder.String(), nil
}
