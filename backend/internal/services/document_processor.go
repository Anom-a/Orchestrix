package services

import (
	"log"
)

func ProcessDocumentInBackground(userID uint, documentID uint, filePath string) {
	log.Printf("Starting background processing for document %d\n", documentID)

	err := ProcessDocumentAI(documentID, filePath)
	if err != nil {
		log.Printf("AI processing failed for document %d: %v\n", documentID, err)

		updateErr := UpdateDocumentStatus(userID, documentID, "failed")
		if updateErr != nil {
			log.Printf("Failed to update document %d status to failed: %v\n", documentID, updateErr)
		}
		return
	}

	err = UpdateDocumentStatus(userID, documentID, "ready")
	if err != nil {
		log.Printf("Failed to update document %d status to ready: %v\n", documentID, err)
		return
	}

	log.Printf("Document %d processed successfully\n", documentID)
}