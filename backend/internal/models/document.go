package models

import "gorm.io/gorm"

type Document struct{
	gorm.Model
	UserID uint `json:"user_id"`
	FileName string `json:"filename"`
	FilePath string `json:"filepath"`
	Status string `json:"status"`
}