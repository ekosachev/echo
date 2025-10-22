package models

import (
	"time"

	"github.com/google/uuid"
)

type EventRrule struct {
	EventID uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Event   Event     `gorm:"foreignKey:EventID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Rrule   string    `gorm:"type:text;not null"`
	RDate   time.Time `gorm:"type:timestamptz[]"`
	ExDate  time.Time `gorm:"type:timestamptz[]"`
	Until   time.Time
	Count   int
}
