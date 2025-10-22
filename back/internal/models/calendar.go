package models

import (
	"time"

	"github.com/google/uuid"
)

type Calendar struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OwnerID   uuid.UUID `gorm:"type:uuid;not null;index:idx_calendars_owner"`
	Owner     User      `gorm:"foreignKey:OwnerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name      string    `gorm:"type:text;not null"`
	Color     string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
