package models

import (
	"time"

	"github.com/google/uuid"
)

type BotLink struct {
	UserID         uuid.UUID `gorm:"primaryKey"`
	TelegramChatID int64     `gorm:"primaryKey;type:bigint"`
	Verified       bool      `gorm:"type:boolean;default:false"`
	LinkCode       string    `gorm:"type:text"`
	ExpiresAt      time.Time `gorm:"type:timestamp"`
}

func (BotLink) TableName() string {
	return "bot_links"
}
