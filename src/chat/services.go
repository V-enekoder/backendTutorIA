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
