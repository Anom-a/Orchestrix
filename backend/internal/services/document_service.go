package services

import (
	"github.com/Anom-a/Orchestrix/internal/database"
	"github.com/Anom-a/Orchestrix/internal/models"
)

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
