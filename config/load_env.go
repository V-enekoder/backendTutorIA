package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error al cargar el archivo .env: %s. Asegúrate de que exista en la raíz del proyecto.", err)
	}
	log.Println(".env cargado correctamente")
}
