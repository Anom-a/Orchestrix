package services

import (
	"errors"
	"time"

	"github.com/Anom-a/Orchestrix/internal/database"
	"github.com/Anom-a/Orchestrix/internal/models"
)

var ErrDocumentNotFound = errors.New("document not found")
var ErrDocumentNotReady = errors.New("document is not ready for querying")
var ErrDocumentFailed = errors.New("document processing failed")

func CreateDocument(userID uint, filename, filepath string) (models.Document, error) {
	doc := models.Document{
		UserID:   userID,
		FileName: filename,
		Filepath: filepath,
		Status:   "uploaded",
	}
	err := database.DB.Create(&doc).Error
	return doc, err
}

func GetUserDocuments(userID uint) ([]models.Document, error) {
	var docs []models.Document
	err := database.DB.Where("user_id = ?", userID).Find(&docs).Error
	return docs, err
}

func GetUserDocumentByID(userID uint, documentID uint) (*models.Document, error) {
	var doc models.Document
	err := database.DB.Where("id = ? AND user_id = ?", documentID, userID).First(&doc).Error
	if err != nil {
		return nil, ErrDocumentNotFound
	}
	if doc.Status == "failed" {
		return nil, ErrDocumentFailed
	}
	if doc.Status != "ready" {
		return nil, ErrDocumentNotReady
	}
	return &doc, nil
}

// GetOwnedDocumentByID returns a document for a user without gating by processing status.
func GetOwnedDocumentByID(userID uint, documentID uint) (*models.Document, error) {
	var doc models.Document
	err := database.DB.Where("id = ? AND user_id = ?", documentID, userID).First(&doc).Error
	if err != nil {
		return nil, ErrDocumentNotFound
	}
	return &doc, nil
}

func UpdateDocumentStatus(userID uint, documentID uint, status string) error {
	return database.DB.Model(&models.Document{}).
		Where("id = ? AND user_id = ?", documentID, userID).
		Update("status", status).Error
}

func MarkDocumentProcessing(userID uint, documentID uint) error {
	now := time.Now()

	return database.DB.Model(&models.Document{}).
		Where("id = ? AND user_id = ?", documentID, userID).
		Updates(map[string]interface{}{
			"status":                "processing",
			"processing_started_at": &now,
			"processed_at":          nil,
		}).Error
}

func MarkDocumentReady(userID uint, documentID uint) error {
	now := time.Now()

	return database.DB.Model(&models.Document{}).
		Where("id = ? AND user_id = ?", documentID, userID).
		Updates(map[string]interface{}{
			"status":       "ready",
			"processed_at": &now,
		}).Error
}

func MarkDocumentFailed(userID uint, documentID uint) error {
	return database.DB.Model(&models.Document{}).
		Where("id = ? AND user_id = ?", documentID, userID).
		Update("status", "failed").Error
}

func FindStaleProcessingDocuments(timeout time.Duration) ([]models.Document, error) {
	var docs []models.Document
	cutoff := time.Now().Add(-timeout)

	err := database.DB.
		Where("status = ? AND processing_started_at IS NOT NULL AND processing_started_at < ?", "processing", cutoff).
		Find(&docs).Error

	return docs, err
}

func CanProcessDocument(status string) bool {
	return status == "uploaded" || status == "failed"
}
