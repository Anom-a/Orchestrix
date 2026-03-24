package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Anom-a/Orchestrix/internal/services"
	"github.com/gin-gonic/gin"
)

type DocumentQueryRequest struct {
	Question string `json:"question"`
}

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

	if statusErr := services.UpdateDocumentStatus(userID, doc.ID, "processing"); statusErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark document as processing"})
		return
	}

	// Call the external AI service
	go func() {
		// Get absolute path
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			fmt.Printf("Failed to get absolute path for %s: %v\n", filePath, err)
			return
		}
		// Convert doc.ID (uint) to string
		docIDStr := fmt.Sprintf("%d", doc.ID)
		if aiErr := services.SendToAIService(docIDStr, absPath); aiErr != nil {
			_ = services.UpdateDocumentStatus(userID, doc.ID, "failed")
			fmt.Printf("Failed to notify AI service for doc %d: %v\n", doc.ID, aiErr)
			return
		}

		if statusErr := services.UpdateDocumentStatus(userID, doc.ID, "ready"); statusErr != nil {
			fmt.Printf("Failed to update status for doc %d: %v\n", doc.ID, statusErr)
		}
	}()

	c.JSON(http.StatusCreated, gin.H{
		"message":     "document uploaded",
		"document_id": doc.ID,
		"file":        doc.FileName,
		"status":      doc.Status,
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

func QueryDocument(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	idParam := c.Param("id")
	documentID64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document id"})
		return
	}
	documentID := uint(documentID64)

	var req DocumentQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Question == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "question is required"})
		return
	}

	doc, err := services.GetUserDocumentByID(userID, documentID)
	if err != nil {
		if errors.Is(err, services.ErrDocumentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		}
		if errors.Is(err, services.ErrDocumentNotReady) {
			c.JSON(http.StatusConflict, gin.H{"error": "document is still processing"})
			return
		}
		if errors.Is(err, services.ErrDocumentFailed) {
			c.JSON(http.StatusConflict, gin.H{"error": "document processing failed"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch document"})
		return
	}

	_ = doc
	if doc.Status != "ready"{
		c.JSON(http.StatusBadRequest, gin.H{"error": "document is not ready for querying"})
	}
	aiResp, err := services.QueryAIService(documentID, req.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query ai service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"document_id": documentID,
		"answer":      aiResp.Answer,
		"sources":     aiResp.Sources,
	})
}

func ListDocumentsByUserID(c *gin.Context){
	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	id := uint(value)
	userID := c.MustGet("user_id").(uint)
	doc, err := services.GetUserDocumentByID(userID, id)
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error getting user from the database" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, doc)
}
func ProcessDocument(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	idParam := c.Param("id")
	documentID64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document id"})
		return
	}
	documentID := uint(documentID64)

	// 1. Fetch document
	doc, err := services.GetUserDocumentByID(userID, documentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	// 2. Validate status
	if doc.Status == "processing" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "document is already processing"})
		return
	}

	if doc.Status == "ready" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "document already processed"})
		return
	}

	// Only allow uploaded or failed
	if doc.Status != "uploaded" && doc.Status != "failed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document state"})
		return
	}

	// 3. Set status = processing
	err = services.UpdateDocumentStatus(userID, documentID, "processing")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	// 4. Call AI service
	err = services.ProcessDocumentAI(documentID, doc.FilePath)
	if err != nil {
		// 5. Mark as failed
		_ = services.UpdateDocumentStatus(userID, documentID, "failed")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "processing failed",
		})
		return
	}

	// 6. Mark as ready
	err = services.UpdateDocumentStatus(userID, documentID, "ready")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update final status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "document processed successfully",
	})
}