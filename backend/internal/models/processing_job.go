package models

import (
	"time"
)

type ProcessingJob struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	DocumentID uint       `gorm:"not null" json:"document_id"`
	Status     string     `gorm:"not null;default:'pending'" json:"status"` // pending, running, completed, failed
	StartedAt  *time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
