package user

import (
	"errors"
	"fmt"

	"github.com/V-enekoder/backendTutorIA/config"
	"github.com/V-enekoder/backendTutorIA/src/schema"
	"gorm.io/gorm"
)

func CreateUserRepository(user schema.User) (uint, error) {
	db := config.DB
	if err := db.Create(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func GetUserByIdRepository(id uint) (schema.User, error) {
	db := config.DB
	var user schema.User

	err := db.Preload("Documents").Preload("Projects").Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schema.User{}, errors.New("record not found")
		}
		return schema.User{}, err
	}

	return user, nil
}

func GetPasswordUserRepository(id uint) (string, error) {
	db := config.DB
	var dbPassword string

	if err := db.Where("id = ?", id).First(&schema.User{}).Error; err != nil {
		return "", err
	}

	if err := db.Model(&schema.User{}).Where("id = ?", id).Pluck("password", &dbPassword).Error; err != nil {
		return "", err
	}
	return dbPassword, nil
}

func UserExistsByFieldRepository(field string, value interface{}, excludeId uint) (bool, error) {
	db := config.DB
	var count int64
	query := fmt.Sprintf("%s = ? AND id != ? AND deleted_at IS NULL", field)
	if err := db.Model(&schema.User{}).Where(query, value, excludeId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func UpdatePasswordUserRepository(id uint, newPassword string) error {
	db := config.DB
	if err := db.Model(&schema.User{}).Where("id = ?", id).Update("password", newPassword).Error; err != nil {
		return err
	}
	return nil
}

func UpdateUserRepository(id uint, userDTO UserUpdateDTO) error {
	db := config.DB

	user := schema.User{}

	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}

	if userDTO.Name != "" {
		user.Name = userDTO.Name
	}
	if userDTO.Email != "" {
		user.Email = userDTO.Email
	}

	if err := db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func DeleteUserbyIDRepository(id uint) error {
	db := config.DB
	var user schema.User

	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}

	if err := db.Where("id = ?", id).Delete(&user).Error; err != nil {
		return err
	}

	return nil
}
