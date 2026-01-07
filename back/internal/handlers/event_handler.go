package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ekosachev/go-backend-template/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type EventHandler struct {
	eventService *service.EventService
}

func NewEventHandler(eventService *service.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
}

type CreateEventRequest struct {
	CalendarID  uuid.UUID `json:"calendar_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartAt     time.Time `json:"start_at" binding:"required"`
	EndAt       time.Time `json:"end_at" binding:"required"`
	Timezone    string    `json:"timezone" binding:"required"`
	AllDay      bool      `json:"all_day"`
	Priority    int       `json:"priority"`
	Color       string    `json:"color"`
}

// CreateEvent godoc
// @Summary Create new event
// @Description Create a new event in calendar with timezone support
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateEventRequest true "Create event request"
// @Success 201 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 403 {object} map[string]any
// @Router /api/v1/events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Errorf("invalid user: %w", err).Error()})
		return
	}

	event, err := h.eventService.CreateEvent(c.Request.Context(), userID, &service.CreateEventRequest{
		CalendarID:  req.CalendarID,
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		StartAt:     req.StartAt,
		EndAt:       req.EndAt,
		Timezone:    req.Timezone,
		AllDay:      req.AllDay,
		Priority:    req.Priority,
		Color:       req.Color,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          event.Id,
		"calendar_id": event.CalendarId,
		"title":       event.Title,
		"description": event.Description,
		"location":    event.Location,
		"start_at":    event.StartAt,
		"end_at":      event.EndAt,
		"timezone":    event.Timezone,
		"all_day":     event.AllDay,
		"priority":    event.Priority,
		"color":       event.Color,
		"created_at":  event.CreatedAt,
	})
}

// GetEvent godoc
// @Summary Get event
// @Description Get event by ID with user's timezone
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Success 200 {object} map[string]any
// @Failure 403 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Router /api/v1/events/{id} [get]
func (h *EventHandler) GetEvent(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	event, err := h.eventService.GetEventForUser(c.Request.Context(), eventID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          event.Id,
		"calendar_id": event.CalendarId,
		"title":       event.Title,
		"description": event.Description,
		"location":    event.Location,
		"start_at":    event.StartAt,
		"end_at":      event.EndAt,
		"timezone":    event.Timezone,
		"all_day":     event.AllDay,
		"priority":    event.Priority,
		"color":       event.Color,
		"created_at":  event.CreatedAt,
		"updated_at":  event.UpdatedAt,
	})
}

type UpdateEventRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartAt     time.Time `json:"start_at" binding:"required"`
	EndAt       time.Time `json:"end_at" binding:"required"`
	Timezone    string    `json:"timezone" binding:"required"`
	AllDay      bool      `json:"all_day"`
	Priority    int       `json:"priority"`
	Color       string    `json:"color"`
}

// UpdateEvent godoc
// @Summary Update event
// @Description Update event with timezone support
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Param request body UpdateEventRequest true "Update event request"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 403 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Router /api/v1/events/{id} [put]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	var req UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	existingEvent, err := h.eventService.GetEventByID(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	event, err := h.eventService.UpdateEventWithTimezone(c.Request.Context(), userID, eventID, &service.CreateEventRequest{
		CalendarID:  existingEvent.CalendarId,
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		StartAt:     req.StartAt,
		EndAt:       req.EndAt,
		Timezone:    req.Timezone,
		AllDay:      req.AllDay,
		Priority:    req.Priority,
		Color:       req.Color,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          event.Id,
		"calendar_id": event.CalendarId,
		"title":       event.Title,
		"description": event.Description,
		"location":    event.Location,
		"start_at":    event.StartAt,
		"end_at":      event.EndAt,
		"timezone":    event.Timezone,
		"all_day":     event.AllDay,
		"priority":    event.Priority,
		"color":       event.Color,
		"updated_at":  event.UpdatedAt,
	})
}

// DeleteEvent godoc
// @Summary Delete event
// @Description Delete event by ID
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Success 204
// @Failure 403 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Router /api/v1/events/{id} [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	_, err = h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	event, err := h.eventService.GetEventByID(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	_ = event

	if err := h.eventService.DeleteEvent(c.Request.Context(), eventID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *EventHandler) getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	claimsAny, exists := c.Get("claims")
	if !exists {
		return uuid.Nil, fmt.Errorf("claims not found")
	}

	claims, ok := claimsAny.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid claims format")
	}

	sub, ok := claims["sub"]
	if !ok {
		return uuid.Nil, fmt.Errorf("sub claim not found")
	}

	userIDStr, ok := sub.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("sub is not a string")
	}

	return uuid.Parse(userIDStr)
}
