package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/V-enekoder/backendTutorIA/src/schema"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dbHost := os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	if dbHost == "" || dbPortStr == "" || dbUser == "" || dbName == "" {
		log.Fatal("Error: Faltan variables de entorno esenciales para la conexión a PostgreSQL (DB_HOST, DB_PORT, DB_USER, DB_NAME)")
	}

	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Error: DB_PORT inválido, debe ser un número: %v", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPort,
		dbSSLMode,
	)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	var gormErr error
	DB, gormErr = gorm.Open(postgres.Open(dsn), gormConfig)

	if gormErr != nil {
		log.Fatalf("Error al conectar con la base de datos PostgreSQL: %v", gormErr)
	}

	log.Println("Conectado exitosamente a la base de datos PostgreSQL!")

}

func SyncDB() {
	if DB == nil {
		log.Fatal("Error: La conexión a la base de datos no ha sido inicializada. Llama a ConnectDB() primero.")
		return
	}

	log.Println("Sincronizando modelos con la base de datos...")
	models := []interface{}{
		&schema.User{},
		&schema.Document{},
		&schema.Project{},
	}

	err := DB.AutoMigrate(models...)
	if err != nil {
		log.Fatalf("Error durante la AutoMigración de GORM: %v", err)
	}
	log.Println("Modelos sincronizados exitosamente.")
}
