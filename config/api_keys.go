package config

import (
	"log"
	"os"
)

func GetGeminiAPIKey() string {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: GEMINI_API_KEY no está definida en las variables de entorno o en el archivo .env.")
	}
	return apiKey
}
