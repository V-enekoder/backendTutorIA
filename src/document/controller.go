package document

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateDocumentController(c *gin.Context) {
	var docDTO DocumentCreateDTO
	if err := c.ShouldBind(&docDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := CreateDocumentService(docDTO)
	if err != nil {
		handleDocumentError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "Document created successfully",
	})
}

func GetDocumentByIdController(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	doc, err := GetDocumentByIdService(uint(id))
	if err != nil {
		handleDocumentError(c, err)
		return
	}

	c.JSON(http.StatusOK, doc)
}

func GetDocumentsByUserController(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	docs, err := GetDocumentsByUserService(uint(userID))
	if err != nil {
		handleDocumentError(c, err)
		return
	}

	c.JSON(http.StatusOK, docs)
}

func UpdateDocumentController(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var docDTO DocumentUpdateDTO
	if err := c.ShouldBindJSON(&docDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := UpdateDocumentService(uint(id), docDTO); err != nil {
		handleDocumentError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Document updated successfully",
	})
}

func DeleteDocumentController(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := DeleteDocumentByIdService(uint(id)); err != nil {
		handleDocumentError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Document deleted successfully",
	})
}

func handleDocumentError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
