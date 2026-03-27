package models

import (
	"time"
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	UserID              uint       `json:"user_id"`
	FileName            string     `json:"filename"`
	Filepath            string     `json:"filepath"`
	Status              string     `json:"status"`
	ProcessingStartedAt *time.Time `json:"processing_started_at"`
	ProcessedAt         *time.Time `json:"processed_at"`
}