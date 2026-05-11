package services

import (
	"errors"
	"time"

	"github.com/Anom-a/Orchestrix/internal/database"
	"github.com/Anom-a/Orchestrix/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrNoPendingJobs = errors.New("no pending jobs found")

func CreateJob(documentID uint) (*models.ProcessingJob, error) {
	job := &models.ProcessingJob{
		DocumentID: documentID,
		Status:     "pending",
	}
	err := database.DB.Create(job).Error
	return job, err
}

// ClaimNextPendingJob safely claims the next pending job atomically using
// PostgreSQL row locks (FOR UPDATE SKIP LOCKED).
func ClaimNextPendingJob() (*models.ProcessingJob, error) {
	var job models.ProcessingJob

	// Start a transaction
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Find the oldest pending job and lock it.
		// FOR UPDATE SKIP LOCKED ensures that if another worker holds a lock
		// on a row, this query will skip it and fetch the next available one.
		err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("status = ?", "pending").
			Order("created_at asc").
			First(&job).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNoPendingJobs
			}
			return err
		}

		// 2. Mark the claimed job as running
		now := time.Now()
		job.Status = "running"
		job.StartedAt = &now

		// Save the updated job status within the same transaction so the row lock applies
		if err := tx.Save(&job).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &job, nil
}

func UpdateJobStatus(jobID uint, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if status == "completed" || status == "failed" {
		now := time.Now()
		updates["finished_at"] = &now
	}

	return database.DB.Model(&models.ProcessingJob{}).Where("id = ?", jobID).Updates(updates).Error
}
