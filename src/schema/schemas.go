package schema

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255)"`
	Email    string `gorm:"type:varchar(255);uniqueIndex"`
	Password string `gorm:"type:varchar(255)"`

	Documents []Document `gorm:"foreignKey:UserID"`
	Projects  []Project  `gorm:"foreignKey:UserID"`
}

type Document struct {
	gorm.Model
	UserID   uint
	Name     string `gorm:"type:varchar(255)"`
	Path     string `gorm:"type:varchar(500)"`
	Resume   string `gorm:"type:text"`
	Mimetype string `gorm:"type:varchar(100)"`
	Size     float64

	User     User       `gorm:"foreignKey:UserID"`
	Projects []*Project `gorm:"many2many:document_project;"`
}

type Project struct {
	gorm.Model
	UserID uint
	Name   string `gorm:"type:varchar(255)"`
	User   User   `gorm:"foreignKey:UserID"`

	Documents []*Document `gorm:"many2many:document_project;"`
}
