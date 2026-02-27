package models

import (
	"time"

	"gorm.io/gorm"
)


type User struct {
  gorm.Model
  FullName string 	`json:"name"`
  Email string 		`gorm:"unique" json:"email"`
  Password string 	`json:"-"`
  Created_at time.Time 	`json:"created_at,omitempty"`
}
