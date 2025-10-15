package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex;size:255;not null"`
	Name      string `gorm:"size:255;not null"`
	Password  string `gorm:"size:255;not null"` // bcrypt hash
	CreatedAt time.Time
	UpdatedAt time.Time
}
