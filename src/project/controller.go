package project

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
    ErrUnauthorized  = errors.New("operación no autorizada")
    ErrProjectNotFound = errors.New("proyecto no encontrado")
)

func CreateProjectController(c *gin.Context) {
	var projectDTO ProjectCreateDTO
	if err := c.ShouldBindJSON(&projectDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := CreateProjectService(projectDTO)
	if err != nil {
		status, errorMessage := handleExceptions(err)
		c.JSON(status, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "Proyecto creado exitosamente",
	})
}

func GetProjectByIdController(c *gin.Context) {
	id := c.Param("id")
	projectId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	project, err := GetProjectByIdService(uint(projectId))
	if err != nil {
		status, errorMessage := handleExceptions(err)
		c.JSON(status, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": project,
	})
}

func UpdateProjectController(c *gin.Context) {
	id := c.Param("id")
	projectId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var projectDTO ProjectUpdateDTO
	if err := c.ShouldBindJSON(&projectDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := UpdateProjectService(uint(projectId), projectDTO); err != nil {
		status, errorMessage := handleExceptions(err)
		c.JSON(status, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Proyecto actualizado exitosamente",
	})
}

func DeleteProjectByIdController(c *gin.Context) {
	id := c.Param("id")
	projectId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := DeleteProjectByIdService(uint(projectId)); err != nil {
		status, errorMessage := handleExceptions(err)
		c.JSON(status, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Proyecto eliminado exitosamente",
	})
}

func handleExceptions(err error) (int, string) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return http.StatusNotFound, "Proyecto no encontrado"
	case errors.Is(err, ErrUnauthorized): // Ejemplo: Si solo el dueño puede modificar
		return http.StatusUnauthorized, "No autorizado"
	default:
		return http.StatusInternalServerError, err.Error()
	}
}