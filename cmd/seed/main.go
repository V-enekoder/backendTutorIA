package main

import (
	"fmt"
	"log"

	"github.com/V-enekoder/backendTutorIA/config"
	"github.com/V-enekoder/backendTutorIA/src/schema"
	"github.com/brianvoe/gofakeit/v7"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
	config.SyncDB()
	db := config.DB

	if err := seedDatabase(db); err != nil {
		log.Fatalf("Error seeding database: %v", err)
	} else {
		log.Println("Database seeded successfully")
	}

}

func seedDatabase(db *gorm.DB) error {
	log.Println("Seeding database...")

	hashedPassword, err := hashPassword("12345")
	if err != nil {
		return err
	}

	// Seed Users
	users, err := seedUsers(db, hashedPassword)
	if err != nil {
		log.Printf("Error seeding users: %v", err)
		return err
	}
	log.Println("Users seeded.")

	// Seed Proyectos
	proyectos, err := seedProyectos(db, users)
	if err != nil {
		log.Printf("Error seeding proyectos: %v", err)
		return err
	}
	log.Println("Proyectos seeded.")

	// Seed Documentos y asociarlos a proyectos
	_, err = seedDocumentos(db, users, proyectos)
	if err != nil {
		log.Printf("Error seeding documentos: %v", err)
		return err
	}
	log.Println("Documentos seeded.")

	return nil
}

// --- Funciones de Seeder ---

// Función para hashear contraseñas
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Seeder para Usuarios
func seedUsers(db *gorm.DB, hashedPassword string) ([]schema.User, error) {
	gofakeit.Seed(0)

	usersToCreate := []schema.User{
		{Nombre: gofakeit.Name(), Correo: gofakeit.Email(), Contraseña: hashedPassword},
		{Nombre: gofakeit.Name(), Correo: gofakeit.Email(), Contraseña: hashedPassword},
		{Nombre: "Usuario Ejemplo", Correo: "ejemplo@test.com", Contraseña: hashedPassword},
	}

	createdUsers := []schema.User{}
	for _, user := range usersToCreate {
		var existingUser schema.User
		// Buscamos o creamos por correo
		result := db.Where(schema.User{Correo: user.Correo}).FirstOrCreate(&existingUser, user)
		if result.Error != nil {
			log.Printf("Error creating or finding user %s: %v", user.Correo, result.Error)
			return nil, fmt.Errorf("failed to seed user with email %s: %w", user.Correo, result.Error)
		}
		createdUsers = append(createdUsers, existingUser)
	}

	return createdUsers, nil
}

// Seeder para Proyectos
func seedProyectos(db *gorm.DB, users []schema.User) ([]schema.Proyecto, error) {
	if len(users) == 0 {
		log.Println("No users provided to seed proyectos.")
		return nil, nil
	}
	gofakeit.Seed(0)

	proyectosToCreate := []schema.Proyecto{
		{Nombre: "Proyecto Alpha", UsuarioID: users[0].ID},     // Asignar al primer usuario
		{Nombre: "Investigación Beta", UsuarioID: users[1].ID}, // Asignar al segundo usuario
	}
	if len(users) > 2 {
		proyectosToCreate = append(proyectosToCreate, schema.Proyecto{Nombre: "Planificación Gamma", UsuarioID: users[2].ID})
	}

	createdProyectos := []schema.Proyecto{}
	for _, proyecto := range proyectosToCreate {
		result := db.Create(&proyecto)
		if result.Error != nil {
			log.Printf("Error creating proyecto %s: %v", proyecto.Nombre, result.Error)
			return nil, fmt.Errorf("failed to seed project %s: %w", proyecto.Nombre, result.Error)
		}
		createdProyectos = append(createdProyectos, proyecto)
	}
	return createdProyectos, nil
}

// Seeder para Documentos
func seedDocumentos(db *gorm.DB, users []schema.User, proyectos []schema.Proyecto) ([]schema.Documento, error) {
	if len(users) == 0 {
		log.Println("No users provided to seed documentos.")
		return nil, nil
	}
	gofakeit.Seed(0)

	documentosToCreate := []schema.Documento{
		{
			Nombre:    "prueba.txt",
			Direccion: "src/documents/prueba.txt",
			Resumen:   gofakeit.LoremIpsumSentence(15),
			Mimetype:  "text/plain",
			Peso:      0.5,
			UsuarioID: users[0].ID,
		},
		{
			Nombre:    "informe_final.pdf",
			Direccion: "src/documents/informe_final.pdf",
			Resumen:   gofakeit.LoremIpsumParagraph(2, 5, 10, " "),
			Mimetype:  "application/pdf",
			Peso:      2.3,
			UsuarioID: users[1].ID,
		},
	}

	createdDocumentos := []schema.Documento{}
	for i, doc := range documentosToCreate {
		result := db.Create(&doc)
		if result.Error != nil {
			return nil, fmt.Errorf("Error creating documento %s: %v", doc.Nombre, result.Error)
		}

		if len(proyectos) > 0 {
			if i == 0 && len(proyectos) > 0 {
				if err := db.Model(&proyectos[0]).Association("Documentos").Append(&doc); err != nil {
					log.Printf("Error associating documento %s to proyecto %s: %v", doc.Nombre, proyectos[0].Nombre, err)
				}
			}
			if i == 1 && len(proyectos) > 0 {
				if err := db.Model(&proyectos[0]).Association("Documentos").Append(&doc); err != nil {
					log.Printf("Error associating documento %s to proyecto %s: %v", doc.Nombre, proyectos[0].Nombre, err)
				}
				if len(proyectos) > 1 {
					if err := db.Model(&proyectos[1]).Association("Documentos").Append(&doc); err != nil {
						log.Printf("Error associating documento %s to proyecto %s: %v", doc.Nombre, proyectos[1].Nombre, err)
					}
				}
			}
		}
		createdDocumentos = append(createdDocumentos, doc)
	}

	return createdDocumentos, nil
}
