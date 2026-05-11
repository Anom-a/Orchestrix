package services

import (
	"log"
	"time"

	"github.com/Anom-a/Orchestrix/internal/database"
	"github.com/Anom-a/Orchestrix/internal/models"
)

// StartJobWorker begins the continuous polling for jobs.
// Multiple instances of this function can be run concurrently across one or more 
// nodes safely because of our FOR UPDATE SKIP LOCKED implementation.
func StartJobWorker(workerID int) {
	log.Printf("Worker %d: starting job polling...", workerID)
	
	// Create infinite worker loop
	for {
		// Attempt to claim a job atomically in a transaction
		job, err := ClaimNextPendingJob()
		
		if err != nil {
			if err == ErrNoPendingJobs {
				// No jobs found, sleep for a short duration before checking again
				// (e.g., 5 seconds) to prevent CPU/Database thrashing
				time.Sleep(5 * time.Second)
				continue
			}
			
			// Some actual DB error occurred
			log.Printf("Worker %d: error claiming job: %v", workerID, err)
			time.Sleep(5 * time.Second)
			continue
		}

		// A job was successfully claimed!
		log.Printf("Worker %d: claimed job ID %d for Document ID %d", workerID, job.ID, job.DocumentID)
		
		// Run actual processing logic
		err = processClaimedJob(workerID, job)
		
		if err != nil {
			log.Printf("Worker %d: failed to process job ID %d: %v", workerID, job.ID, err)
			UpdateJobStatus(job.ID, "failed")
		} else {
			log.Printf("Worker %d: successfully completed job ID %d", workerID, job.ID)
			UpdateJobStatus(job.ID, "completed")
		}
	}
}

// processClaimedJob encapsulates the execution mapping a job ID to the document processing layer
func processClaimedJob(workerID int, job *models.ProcessingJob) error {
	// 1. Fetch the actual Document from the database to get filepath and owner
	var doc models.Document
	if err := database.DB.First(&doc, job.DocumentID).Error; err != nil {
		return err // Document not found or db error
	}
	
	// Tie together the document and job systems
	// Update Document status to processing
	err := MarkDocumentProcessing(doc.UserID, doc.ID)
	if err != nil {
		log.Printf("Worker %d: failed to mark doc processing: %v", workerID, err)
		// continue anyway to try the heavy lifting
	}
	
	// 2. Do the heavy lifting (AI processing via Python service)
	log.Printf("Worker %d: processing document %d at path %s", workerID, doc.ID, doc.Filepath)
	
	// ProcessDocumentAI should perform synchronous blocking python calls 
	err = ProcessDocumentAI(doc.ID, doc.Filepath)
	
	if err != nil {
		// Log and update to failed on document side
		MarkDocumentFailed(doc.UserID, doc.ID)
		return err // bubble up so job side becomes failed as well
	}
	
	// 3. Update document status upon finish
	err = MarkDocumentReady(doc.UserID, doc.ID)
	if err != nil {
		log.Printf("Worker %d: failed to mark doc ready: %v", workerID, err)
		return err
	}
	
	return nil
}
