package config

import (
	"fmt"
	"log"
	"os"
	"strconv" // Para convertir DB_PORT a int

	"github.com/V-enekoder/TutorIA/src/schema"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // Opcional, para configurar el logger de GORM
)

var DB *gorm.DB

func ConnectDB() {
	// Cargar variables de entorno para PostgreSQL
	dbHost := os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")
	// dbTimeZone := os.Getenv("DB_TIMEZONE") // Opcional, si la defines

	// Validar que las variables esenciales estén presentes
	if dbHost == "" || dbPortStr == "" || dbUser == "" || dbName == "" {
		log.Fatal("Error: Faltan variables de entorno esenciales para la conexión a PostgreSQL (DB_HOST, DB_PORT, DB_USER, DB_NAME)")
	}

	// Convertir puerto a entero (GORM DSN lo espera así usualmente)
	// Aunque el DSN string puede tomarlo como string, es buena práctica.
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Error: DB_PORT inválido, debe ser un número: %v", err)
	}

	// Construir la cadena de conexión DSN para PostgreSQL
	// Formato: "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dbHost,
		dbUser,
		dbPassword, // La contraseña puede estar vacía si así está configurada la BD para ese usuario/host
		dbName,
		dbPort,
		dbSSLMode,
	)
	// if dbTimeZone != "" { // Añadir TimeZone si está definida
	// 	dsn += " TimeZone=" + dbTimeZone
	// }

	// Configuración de GORM (puedes personalizar el logger aquí)
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Muestra todas las SQL queries. Cambia a logger.Silent en producción si no las necesitas.
	}

	// Abrir la conexión a la base de datos
	var gormErr error
	DB, gormErr = gorm.Open(postgres.Open(dsn), gormConfig)

	if gormErr != nil {
		log.Fatalf("Error al conectar con la base de datos PostgreSQL: %v", gormErr)
	}

	log.Println("Conectado exitosamente a la base de datos PostgreSQL!")

	// (Opcional) Configurar el pool de conexiones si tienes esas variables en .env
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Error al obtener el objeto sql.DB subyacente: %v", err)
	}

	if maxOpenConnsStr := os.Getenv("DB_MAX_OPEN_CONNS"); maxOpenConnsStr != "" {
		maxOpenConns, _ := strconv.Atoi(maxOpenConnsStr)
		sqlDB.SetMaxOpenConns(maxOpenConns)
	}
	if maxIdleConnsStr := os.Getenv("DB_MAX_IDLE_CONNS"); maxIdleConnsStr != "" {
		maxIdleConns, _ := strconv.Atoi(maxIdleConnsStr)
		sqlDB.SetMaxIdleConns(maxIdleConns)
	}
	// ... puedes añadir SetConnMaxIdleTime y SetConnMaxLifetime de manera similar
}

func SyncDB() {
	if DB == nil {
		log.Fatal("Error: La conexión a la base de datos no ha sido inicializada. Llama a ConnectDB() primero.")
		return
	}

	log.Println("Sincronizando modelos con la base de datos...")
	models := []interface{}{
		&schema.User{},
		&schema.Review{},
		&schema.ReviewImage{},
		&schema.Place{},
		&schema.Comment{},
		&schema.Answer{},
		&schema.Reaction{},
		&schema.Notification{},
		&schema.ValidationCode{},
	}

	err := DB.AutoMigrate(models...)
	if err != nil {
		log.Fatalf("Error durante la AutoMigración de GORM: %v", err)
	}
	log.Println("Modelos sincronizados exitosamente.")
}

// (Opcional) Función para cargar la variable GEMINI_API_KEY si la necesitas en este paquete
// o podrías tener un paquete `services` o `utils` para esto.
func GetGeminiAPIKey() string {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("Advertencia: GEMINI_API_KEY no está definida en las variables de entorno.")
	}
	return apiKey
}
