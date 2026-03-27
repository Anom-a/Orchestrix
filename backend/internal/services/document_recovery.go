package services

import (
	"log"
	"time"
)

func RecoverStaleProcessingDocuments(timeout time.Duration) (int, error) {
	staleDocs, err := FindStaleProcessingDocuments(timeout)
	if err != nil {
		return 0, err
	}

	recovered := 0
	for _, doc := range staleDocs {
		if err := MarkDocumentFailed(doc.UserID, doc.ID); err != nil {
			log.Printf("Failed to mark stale document %d as failed: %v", doc.ID, err)
			continue
		}
		recovered++
	}

	return recovered, nil
}
