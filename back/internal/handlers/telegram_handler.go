package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ekosachev/go-backend-template/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TelegramHandler struct {
	db *gorm.DB
}

func NewTelegramHandler(db *gorm.DB) *TelegramHandler {
	return &TelegramHandler{
		db: db,
	}
}

type LinkTelegramRequest struct {
	TelegramChatID int64 `json:"telegram_chat_id" binding:"required"`
}

func (h *TelegramHandler) LinkTelegram(c *gin.Context) {
	var req LinkTelegramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
	}
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user",
		})
	}

	botLink := &models.BotLink{
		UserID:         userID,
		TelegramChatID: req.TelegramChatID,
		Verified:       true,
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	if err := h.db.Create(botLink).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to link telegram"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "telegram linked successfully",
		"telegram_chat_id": req.TelegramChatID,
	})

}

func (h *TelegramHandler) getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	claimsAny, exists := c.Get("claims")
	if !exists {
		return uuid.Nil, fmt.Errorf("claims not found")
	}
	claims, ok := claimsAny.(map[string]any)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid claims")
	}
	userIDStr, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("user ID not found")
	}
	return uuid.Parse(userIDStr)
}
