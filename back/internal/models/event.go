package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Event struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CalendarID  uuid.UUID      `gorm:"type:uuid;not null;index:idx_events_calendar"`
	Calendar    Calendar       `gorm:"foreignKey:CalendarID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Title       string         `gorm:"type:text;not null"`
	Description string         `gorm:"type:text"`
	Location    string         `gorm:"type:text"`
	StartAt     time.Time      `gorm:"not null;index:idx_events_start_end"`
	EndAt       time.Time      `gorm:"not null;index:idx_events_start_end"`
	Timezone    string         `gorm:"type:text;default:MSK"`
	AllDay      bool           `gorm:"default:false"`
	Priority    int            `gorm:"default:0"`
	Color       string         `gorm:"type:text"`
	Meta        datatypes.JSON `gorm:"type:jsonb;default:'{}'::jsonb"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
}
