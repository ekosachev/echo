package models

import (
	"time"

	"github.com/google/uuid"
)

type EventRrule struct {
	EventId uuid.UUID   `gorm:"type:uuid;primaryKey"`
	Rrule   string      `gorm:"type:text"`
	Rdate   []time.Time `gorm:"type:timestamp[]"`
	Exdate  []time.Time `gorm:"type:timestamp[]"`
	Until   time.Time   `gorm:"type:timestamp"`
	Count   int         `gorm:"type:int"`
}
