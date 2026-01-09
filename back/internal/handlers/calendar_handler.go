package handlers

import (
	"fmt"
	"net/http"

	"github.com/ekosachev/go-backend-template/internal/models"
	"github.com/ekosachev/go-backend-template/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CalendarHandler struct {
	calendarService *service.CalendarService
	eventService    *service.EventService
}

func NewCalendarHandler(
	calendarService *service.CalendarService,
	eventService *service.EventService,
) *CalendarHandler {
	return &CalendarHandler{
		calendarService: calendarService,
		eventService:    eventService,
	}
}

type CreateCalendarRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=255"`
	Color string `json:"color"`
}

// CreateCalendar godoc
// @Summary Create new calendar
// @Description Create a new calendar for authenticated user
// @Tags calendars
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCalendarRequest true "Create calendar request"
// @Success 201 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/v1/calendars [post]
func (h *CalendarHandler) CreateCalendar(c *gin.Context) {
	var req CreateCalendarRequest
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

	calendar := &models.Calendar{
		OwnerID: userID,
		Name:    req.Name,
		Color:   req.Color,
	}

	if err := h.calendarService.CreateCalendar(c.Request.Context(), calendar); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       calendar.ID,
		"name":     calendar.Name,
		"color":    calendar.Color,
		"owner_id": calendar.OwnerID,
	})
}

type UpdateCalendarRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=255"`
	Color string `json:"color"`
}

// UpdateCalendar godoc
// @Summary Update calendar
// @Description Update calendar information
// @Tags calendars
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Calendar ID"
// @Param request body UpdateCalendarRequest true "Update calendar request"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 403 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Router /api/v1/calendars/{id} [put]
func (h *CalendarHandler) UpdateCalendar(c *gin.Context) {
	calendarID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid calendar id"})
		return
	}

	var req UpdateCalendarRequest
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

	calendar, err := h.calendarService.GetCalendarForUser(c.Request.Context(), calendarID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "calendar not found"})
		return
	}

	calendar.Name = req.Name
	calendar.Color = req.Color

	if err := h.calendarService.UpdateCalendar(c.Request.Context(), calendar); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       calendar.ID,
		"name":     calendar.Name,
		"color":    calendar.Color,
		"owner_id": calendar.OwnerID,
	})
}

// GetCalendar godoc
// @Summary Get calendar
// @Description Get calendar by ID
// @Tags calendars
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Calendar ID"
// @Success 200 {object} map[string]any
// @Failure 403 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Router /api/v1/calendars/{id} [get]
func (h *CalendarHandler) GetCalendar(c *gin.Context) {
	calendarID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid calendar id"})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	calendar, err := h.calendarService.GetCalendarForUser(c.Request.Context(), calendarID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Errorf("calendar not found: %w", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         calendar.ID,
		"name":       calendar.Name,
		"color":      calendar.Color,
		"owner_id":   calendar.OwnerID,
		"created_at": calendar.CreatedAt,
	})
}

// GetUserCalendars godoc
// @Summary Get user calendars
// @Description Get all calendars for authenticated user
// @Tags calendars
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/calendars [get]
func (h *CalendarHandler) GetUserCalendars(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	calendars, err := h.calendarService.GetUserCalendars(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]map[string]any, len(calendars))
	for i, calendar := range calendars {
		response[i] = map[string]any{
			"id":         calendar.ID,
			"name":       calendar.Name,
			"color":      calendar.Color,
			"owner_id":   calendar.OwnerID,
			"created_at": calendar.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"calendars": response})
}

// DeleteCalendar godoc
// @Summary Delete calendar
// @Description Delete calendar and all its events
// @Tags calendars
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Calendar ID"
// @Success 204
// @Failure 403 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Router /api/v1/calendars/{id} [delete]
func (h *CalendarHandler) DeleteCalendar(c *gin.Context) {
	calendarID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid calendar id"})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	_, err = h.calendarService.GetCalendarForUser(c.Request.Context(), calendarID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "calendar not found"})
		return
	}

	if err := h.calendarService.DeleteCalendar(c.Request.Context(), calendarID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetCalendarEvents godoc
// @Summary Get calendar events
// @Description Get all events for specific calendar
// @Tags calendars, events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Calendar ID"
// @Success 200 {object} map[string]any
// @Failure 403 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Router /api/v1/calendars/{id}/events [get]
func (h *CalendarHandler) GetCalendarEvents(c *gin.Context) {
	calendarID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid calendar id"})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	_, err = h.calendarService.GetCalendarForUser(c.Request.Context(), calendarID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "calendar not found"})
		return
	}

	events, err := h.eventService.GetEventsByCalendar(c.Request.Context(), calendarID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]map[string]any, len(events))
	for i, event := range events {
		response[i] = map[string]any{
			"id":          event.Id,
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
		}
	}

	c.JSON(http.StatusOK, gin.H{"events": response})
}

func (h *CalendarHandler) getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
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
