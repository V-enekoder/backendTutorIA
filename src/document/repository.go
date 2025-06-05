package document

import (
	"errors"
	"fmt"

	"github.com/V-enekoder/backendTutorIA/config"
	"github.com/V-enekoder/backendTutorIA/src/schema"
	"gorm.io/gorm"
)

func CreateDocumentRepository(doc schema.Document) (uint, error) {
	db := config.DB

	if err := db.Create(&doc).Error; err != nil {
		return 0, fmt.Errorf("error al crear documento: %w", err)
	}
	return doc.ID, nil
}

// obtiene un documento por su ID con relaciones cargadas
func GetDocumentByIdRepository(id uint) (schema.Document, error) {
	db := config.DB
	var doc schema.Document

	err := db.Preload("User").Where("id = ?", id).First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schema.Document{}, errors.New("documento no encontrado")
		}
		return schema.Document{}, fmt.Errorf("error al obtener documento: %w", err)
	}
	return doc, nil
}

// Verifica si existe un documento con un valor específico en un campo
func DocumentExistsByFieldRepository(field string, value interface{}, excludeId uint) (bool, error) {
	db := config.DB
	var count int64

	query := fmt.Sprintf("%s = ? AND id != ? AND deleted_at IS NULL", field)
	if err := db.Model(&schema.Document{}).Where(query, value, excludeId).Count(&count).Error; err != nil {
		return false, fmt.Errorf("error al verificar existencia de documento: %w", err)
	}
	return count > 0, nil
}

func UpdateDocumentRepository(id uint, docDTO DocumentUpdateDTO) error {
	db := config.DB

	return db.Transaction(func(tx *gorm.DB) error {
		var doc schema.Document
		if err := tx.Where("id = ?", id).First(&doc).Error; err != nil {
			return fmt.Errorf("documento no encontrado: %w", err)
		}

		// Actualizar solo campos proporcionados
		if docDTO.Name != "" {
			doc.Name = docDTO.Name
		}
		if docDTO.Address != "" {
			doc.Path = docDTO.Address
		}
		if docDTO.Resume != "" {
			doc.Resume = docDTO.Resume
		}
		if docDTO.Size != 0 {
			doc.Size = docDTO.Size
		}

		if err := tx.Save(&doc).Error; err != nil {
			return fmt.Errorf("error al guardar cambios: %w", err)
		}

		return nil
	})
}

func DeleteDocumentByIdRepository(id uint) error {
	db := config.DB
	var doc schema.Document

	// verificar que el documento existe
	if err := db.Where("id = ?", id).First(&doc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("documento no encontrado")
		}
		return fmt.Errorf("error al buscar documento: %w", err)
	}

	// Borrado físico (para borrado lógico usar db.Delete(&doc))
	if err := db.Unscoped().Delete(&doc).Error; err != nil {
		return fmt.Errorf("error al eliminar documento: %w", err)
	}

	return nil
}

// obtiene todos los documentos de un usuario
func GetDocumentsByUserRepository(userID uint) ([]schema.Document, error) {
	db := config.DB
	var documents []schema.Document

	err := db.Where("user_id = ?", userID).Find(&documents).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener documentos del usuario: %w", err)
	}
	return documents, nil
}

// obtiene documentos asociados a un proyecto
func GetDocumentsByProjectRepository(projectID uint) ([]schema.Document, error) {
	db := config.DB
	var documents []schema.Document

	err := db.Where("project_id = ?", projectID).Find(&documents).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener documentos del proyecto: %w", err)
	}
	return documents, nil
}
