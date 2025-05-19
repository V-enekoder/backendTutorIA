package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("Debes crear un archivo .env en la raíz del proyecto")
	}
}
