package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateUserController(c *gin.Context) {
	var userDTO UserCreateDTO
	if err := c.ShouldBind(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := CreateUserService(userDTO)
	if err != nil {
		status, errorMessage := handleExceptions(err)
		c.JSON(status, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})

}

func handleExceptions(err error) (int, string) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return http.StatusNotFound, "Usuario no Encontrado"
	default:
		return http.StatusInternalServerError, err.Error()
	}
}

func GetUserByIdController(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, err := GetUserByIdService(uint(userId))
	if err != nil {
		status, errorMessage := handleExceptions(err)
		c.JSON(status, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func UpdatePasswordUserController(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var passwordDTO UserUpdatePasswordDTO
	if err := c.ShouldBindJSON(&passwordDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := UpdatePasswordUserService(uint(userId), passwordDTO); err != nil {
		status, errorMessage := handleExceptions(err)
		c.JSON(status, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

func UpdateUserController(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var userDTO UserUpdateDTO
	if err := c.ShouldBind(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := UpdateUserService(uint(userId), userDTO); err != nil {
		status, errorMessage := handleExceptions(err)
		c.JSON(status, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
	})
}

func DeleteUserbyIdController(c *gin.Context) {
	id := c.Param("id")
	revId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = DeleteUserByIdService(uint(revId))
	if err != nil {
		status, errorMessage := handleExceptions(err)
		c.JSON(status, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Usuario eliminado Exitosamente",
	})
}
