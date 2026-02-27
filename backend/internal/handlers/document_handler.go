package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/Anom-a/Orchestrix/internal/services"
	"github.com/gin-gonic/gin"
)

func UploadDocument(c *gin.Context){
	userID := c.MustGet("user_id").(uint)
	file, err := c.FormFile("file")
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	uploadDir := "storage/uploads"
	os.MkdirAll(uploadDir, os.ModePerm)
	filePath := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, filePath);err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create the document"})
		return
	}
	err = services.CreateDocument(userID, file.Filename, filePath)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save the document to the db"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "document uploaded",
		"file": file.Filename,
	})
}

func ListDocuments(c *gin.Context){
	userID := c.MustGet("user_id").(uint)
	docs, err := services.GetUserDocuments(userID)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch documents"})
		return
	}
	c.JSON(http.StatusOK, docs)
}