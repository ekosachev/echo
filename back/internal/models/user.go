package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email        string    `gorm:"type:text;index:idx_users_email,unique;not null"`
	PasswordHash string    `gorm:"type:text;not null"` // bcrypt hash
	Timezone     string    `gorm:"type:text;default:MSK"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
