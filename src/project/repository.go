package project

import (
	"errors"
	"fmt"

	"github.com/V-enekoder/backendTutorIA/config"
	"github.com/V-enekoder/backendTutorIA/src/schema"
	"gorm.io/gorm"
)


func CreateProjectRepository(project schema.Project) (uint, error) {
	db := config.DB
	if err := db.Create(&project).Error; err != nil {
		return 0, err
	}
	return project.ID, nil
}


func GetProjectByIdRepository(id uint) (schema.Project, error) {
	db := config.DB
	var project schema.Project

	err := db.Preload("Documents").Where("id = ?", id).First(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schema.Project{}, errors.New("proyecto no encontrado")
		}
		return schema.Project{}, err
	}

	return project, nil
}


func ProjectExistsByFieldRepository(field string, value interface{}, excludeId uint) (bool, error) {
	db := config.DB
	var count int64
	query := fmt.Sprintf("%s = ? AND id != ? AND deleted_at IS NULL", field)
	if err := db.Model(&schema.Project{}).Where(query, value, excludeId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}


func UpdateProjectRepository(id uint, projectDTO ProjectUpdateDTO) error {
	db := config.DB

	project := schema.Project{}
	if err := db.Where("id = ?", id).First(&project).Error; err != nil {
		return err
	}

	if projectDTO.Name != "" {
		project.Name = projectDTO.Name
	}
	if projectDTO.Address != "" {
		project.Address = projectDTO.Address
	}
	if projectDTO.Summary != "" {
		project.Summary = projectDTO.Summary
	}
	if projectDTO.MimeType != "" {
		project.MimeType = projectDTO.MimeType
	}
	if projectDTO.FileSize != 0 {
		project.FileSize = projectDTO.FileSize
	}

	if err := db.Save(&project).Error; err != nil {
		return err
	}

	return nil
}


func DeleteProjectByIdRepository(id uint) error {
	db := config.DB
	var project schema.Project

	if err := db.Where("id = ?", id).First(&project).Error; err != nil {
		return err
	}

	// Borrado físico (opcional: usar db.Delete(&project) para borrado lógico)
	if err := db.Unscoped().Delete(&project).Error; err != nil {
		return err
	}

	return nil
}