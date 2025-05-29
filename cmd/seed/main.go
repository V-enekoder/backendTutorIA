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
	proyectos, err := seedProjects(db, users)
	if err != nil {
		log.Printf("Error seeding proyectos: %v", err)
		return err
	}
	log.Println("Proyectos seeded.")

	// Seed Documentos y asociarlos a proyectos
	_, err = seedDocuments(db, users, proyectos)
	if err != nil {
		log.Printf("Error seeding documentos: %v", err)
		return err
	}
	log.Println("Documentos seeded.")

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Seeder para Usuarios
func seedUsers(db *gorm.DB, hashedPassword string) ([]schema.User, error) {
	gofakeit.Seed(0)

	usersToCreate := []schema.User{
		{Name: gofakeit.Name(), Email: gofakeit.Email(), Password: hashedPassword},
		{Name: gofakeit.Name(), Email: gofakeit.Email(), Password: hashedPassword},
		{Name: "Usuario Ejemplo", Email: "ejemplo@test.com", Password: hashedPassword},
	}

	createdUsers := []schema.User{}
	for _, user := range usersToCreate {
		var existingUser schema.User
		// Buscamos o creamos por correo
		result := db.Where(schema.User{Email: user.Email}).FirstOrCreate(&existingUser, user)
		if result.Error != nil {
			log.Printf("Error creating or finding user %s: %v", user.Email, result.Error)
			return nil, fmt.Errorf("failed to seed user with email %s: %w", user.Email, result.Error)
		}
		createdUsers = append(createdUsers, existingUser)
	}

	return createdUsers, nil
}

// Seeder para Proyectos
func seedProjects(db *gorm.DB, users []schema.User) ([]schema.Project, error) {
	if len(users) == 0 {
		log.Println("No users provided to seed proyectos.")
		return nil, nil
	}
	gofakeit.Seed(0)

	proyectosToCreate := []schema.Project{
		{Name: "Proyecto Alpha", UserID: users[0].ID},
		{Name: "Investigación Beta", UserID: users[1].ID},
	}
	if len(users) > 2 {
		proyectosToCreate = append(proyectosToCreate, schema.Project{
			Name:   "Planificación Gamma",
			UserID: users[2].ID})
	}

	createdProyectos := []schema.Project{}
	for _, proyecto := range proyectosToCreate {
		result := db.Create(&proyecto)
		if result.Error != nil {
			log.Printf("Error creating proyecto %s: %v", proyecto.Name, result.Error)
			return nil, fmt.Errorf("failed to seed project %s: %w", proyecto.Name, result.Error)
		}
		createdProyectos = append(createdProyectos, proyecto)
	}
	return createdProyectos, nil
}

// Seeder para Documentos
func seedDocuments(db *gorm.DB, users []schema.User, proyectos []schema.Project) ([]schema.Document, error) {
	if len(users) == 0 {
		log.Println("No users provided to seed documentos.")
		return nil, nil
	}
	gofakeit.Seed(0)

	documentosToCreate := []schema.Document{
		{
			Name:     "prueba.txt",
			Path:     "src/documents/prueba.txt",
			Resume:   gofakeit.LoremIpsumSentence(15),
			Mimetype: "text/plain",
			Size:     0.5,
			UserID:   users[0].ID,
		},
		{
			Name:     "informe_final.pdf",
			Path:     "src/documents/informe_final.pdf",
			Resume:   gofakeit.LoremIpsumParagraph(2, 5, 10, " "),
			Mimetype: "application/pdf",
			Size:     2.3,
			UserID:   users[1].ID,
		},
	}

	createdDocumentos := []schema.Document{}
	for i, doc := range documentosToCreate {
		result := db.Create(&doc)
		if result.Error != nil {
			return nil, fmt.Errorf("Error creating documento %s: %v", doc.Name, result.Error)
		}

		if len(proyectos) > 0 {
			if i == 0 && len(proyectos) > 0 {
				if err := db.Model(&proyectos[0]).Association("Documents").Append(&doc); err != nil {
					log.Printf("Error associating documento %s to proyecto %s: %v", doc.Name, proyectos[0].Name, err)
				}
			}
			if i == 1 && len(proyectos) > 0 {
				if err := db.Model(&proyectos[0]).Association("Documents").Append(&doc); err != nil {
					log.Printf("Error associating documento %s to proyecto %s: %v", doc.Name, proyectos[0].Name, err)
				}
				if len(proyectos) > 1 {
					if err := db.Model(&proyectos[1]).Association("Documents").Append(&doc); err != nil {
						log.Printf("Error associating documento %s to proyecto %s: %v", doc.Name, proyectos[1].Name, err)
					}
				}
			}
		}
		createdDocumentos = append(createdDocumentos, doc)
	}

	return createdDocumentos, nil
}
