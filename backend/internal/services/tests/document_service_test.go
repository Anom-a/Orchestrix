package services

import (
	"errors"
	"testing"

	"github.com/Anom-a/Orchestrix/internal/database"
	"github.com/Anom-a/Orchestrix/internal/models"
	"github.com/Anom-a/Orchestrix/internal/services"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDocumentTestDB(t *testing.T) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite test db: %v", err)
	}

	if err := db.AutoMigrate(&models.Document{}); err != nil {
		t.Fatalf("failed to migrate document model: %v", err)
	}

	database.DB = db
}

func TestCreateDocumentPersistsWithUploadedStatus(t *testing.T) {
	setupDocumentTestDB(t)

	doc, err := services.CreateDocument(7, "notes.txt", "storage/uploads/notes.txt")
	if err != nil {
		t.Fatalf("CreateDocument returned error: %v", err)
	}

	if doc.Status != "uploaded" {
		t.Fatalf("expected status uploaded, got %q", doc.Status)
	}

	if doc.UserID != 7 {
		t.Fatalf("expected user id 7, got %d", doc.UserID)
	}
}

func TestGetUserDocumentsFiltersByUserID(t *testing.T) {
	setupDocumentTestDB(t)

	if _, err := services.CreateDocument(1, "a.txt", "storage/uploads/a.txt"); err != nil {
		t.Fatalf("failed to create document A: %v", err)
	}
	if _, err := services.CreateDocument(1, "b.txt", "storage/uploads/b.txt"); err != nil {
		t.Fatalf("failed to create document B: %v", err)
	}
	if _, err := services.CreateDocument(2, "c.txt", "storage/uploads/c.txt"); err != nil {
		t.Fatalf("failed to create document C: %v", err)
	}

	docs, err := services.GetUserDocuments(1)
	if err != nil {
		t.Fatalf("GetUserDocuments returned error: %v", err)
	}

	if len(docs) != 2 {
		t.Fatalf("expected 2 documents for user 1, got %d", len(docs))
	}

	for _, doc := range docs {
		if doc.UserID != 1 {
			t.Fatalf("found document with unexpected user id: %d", doc.UserID)
		}
	}
}

func TestGetUserDocumentByIDReady(t *testing.T) {
	setupDocumentTestDB(t)

	doc := models.Document{
		UserID:   1,
		FileName: "ready.txt",
		Filepath: "storage/uploads/ready.txt",
		Status:   "ready",
	}
	if err := database.DB.Create(&doc).Error; err != nil {
		t.Fatalf("failed to seed ready doc: %v", err)
	}

	loaded, err := services.GetUserDocumentByID(1, doc.ID)
	if err != nil {
		t.Fatalf("expected ready doc, got error: %v", err)
	}

	if loaded.ID != doc.ID {
		t.Fatalf("expected document id %d, got %d", doc.ID, loaded.ID)
	}
}

func TestGetUserDocumentByIDProcessing(t *testing.T) {
	setupDocumentTestDB(t)

	doc := models.Document{
		UserID:   1,
		FileName: "processing.txt",
		Filepath: "storage/uploads/processing.txt",
		Status:   "processing",
	}
	if err := database.DB.Create(&doc).Error; err != nil {
		t.Fatalf("failed to seed processing doc: %v", err)
	}

	_, err := services.GetUserDocumentByID(1, doc.ID)
	if !errors.Is(err, services.ErrDocumentNotReady) {
		t.Fatalf("expected ErrDocumentNotReady, got %v", err)
	}
}

func TestGetUserDocumentByIDFailed(t *testing.T) {
	setupDocumentTestDB(t)

	doc := models.Document{
		UserID:   1,
		FileName: "failed.txt",
		Filepath: "storage/uploads/failed.txt",
		Status:   "failed",
	}
	if err := database.DB.Create(&doc).Error; err != nil {
		t.Fatalf("failed to seed failed doc: %v", err)
	}

	_, err := services.GetUserDocumentByID(1, doc.ID)
	if !errors.Is(err, services.ErrDocumentFailed) {
		t.Fatalf("expected ErrDocumentFailed, got %v", err)
	}
}
