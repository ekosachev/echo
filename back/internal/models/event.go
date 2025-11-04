package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Id          uuid.UUID              `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CalendarId  uuid.UUID              `gorm:"type:uuid;not null;on delete cascade index:idx_events_calendar"`
	Title       string                 `gorm:"type:text;not null"`
	Description string                 `gorm:"type:text"`
	Location    string                 `gorm:"type:text"`
	StartAt     time.Time              `gorm:"type:timestamp;not null;index:idx_events_start_end"`
	EndAt       time.Time              `gorm:"type:timestamp;not null;index:idx_events_start_end"`
	Timezone    string                 `gorm:"type:text;default:'UTC'"`
	AllDay      bool                   `gorm:"type:boolean;default:false"`
	Priority    int                    `gorm:"type:int;default:0"`
	Color       string                 `gorm:"type:text"`
	Meta        map[string]interface{} `gorm:"type:jsonb;default:'{}'"`
	CreatedAt   time.Time              `gorm:"type:timestamp;default:now()"`
	UpdatedAt   time.Time              `gorm:"type:timestamp;default:now()"`
}
