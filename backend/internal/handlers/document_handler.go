package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Anom-a/Orchestrix/internal/services"
	"github.com/gin-gonic/gin"
)

func UploadDocument(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	uploadDir := "storage/uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating file"})
		return
	}
	uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	filePath := filepath.Join(uploadDir, uniqueName)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save the document"})
		return
	}
	doc, err := services.CreateDocument(userID, uniqueName, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save the create to the db"})
		return
	}

	// Call the external AI service
	go func() {
		// Convert doc.ID (uint) to string
		docIDStr := fmt.Sprintf("%d", doc.ID)
		if aiErr := services.SendToAIService(docIDStr, filePath); aiErr != nil {
			fmt.Printf("Failed to notify AI service for doc %d: %v\n", doc.ID, aiErr)
		}
	}()

	c.JSON(http.StatusCreated, gin.H{
		"message": "document uploaded",
		"file":    file.Filename,
	})
}

func ListDocuments(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	docs, err := services.GetUserDocuments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch documents"})
		return
	}
	c.JSON(http.StatusOK, docs)
}
