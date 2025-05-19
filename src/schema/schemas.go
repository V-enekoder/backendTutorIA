package schema

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Nombre     string `gorm:"type:varchar(255)"`
	Correo     string `gorm:"type:varchar(255);uniqueIndex"`
	Contraseña string `gorm:"type:varchar(255)"`

	Documentos []Documento `gorm:"foreignKey:UsuarioID"`
	Proyectos  []Proyecto  `gorm:"foreignKey:UsuarioID"`
}

type Documento struct {
	gorm.Model
	UsuarioID uint
	Nombre    string `gorm:"type:varchar(255)"`
	Direccion string `gorm:"type:varchar(500)"`
	Resumen   string `gorm:"type:text"`
	Mimetype  string `gorm:"type:varchar(100)"`
	Peso      float64

	Usuario   User        `gorm:"foreignKey:UsuarioID"`
	Proyectos []*Proyecto `gorm:"many2many:documento_proyecto;"`
}

type Proyecto struct {
	gorm.Model
	UsuarioID uint
	Nombre    string `gorm:"type:varchar(255)"`
	Usuario   User   `gorm:"foreignKey:UsuarioID"`

	Documentos []*Documento `gorm:"many2many:documento_proyecto;"`
}
