package services

import (
	"errors"

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
		FilePath: filepath,
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
func CanProcessDocument(status string) bool {
	return status == "uploaded" || status == "failed"
}
