package handlers

import (
	"net/http"

	"github.com/ekosachev/go-backend-template/internal/models"
	"github.com/ekosachev/go-backend-template/internal/service"
	"github.com/gin-gonic/gin"
)

type CalendarHandler struct {
	service *service.CalendarService
}

func NewCalendarHandler(service *service.CalendarService) *CalendarHandler {
	return &CalendarHandler{service: service}
}

func (h *CalendarHandler) CreateCalendar(ctx *gin.Context) {
	var calendar models.Calendar

	if err := ctx.ShouldBindJSON(&calendar); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	if err := h.service.CreateCalendar(ctx.Request.Context(), &calendar); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, calendar)
}

func (h *CalendarHandler) GetCalendarByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "calendar ID is empty"})
		return
	}

	calendar, err := h.service.GetCalendarByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, calendar)
}

func (h *CalendarHandler) UpdateCalendar(ctx *gin.Context) {
	var calendar models.Calendar

	if err := ctx.ShouldBindJSON(&calendar); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	if err := h.service.UpdateCalendar(ctx.Request.Context(), &calendar); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, calendar)
}

func (h *CalendarHandler) DeleteCalendar(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "calendar ID is empty"})
		return
	}

	if err := h.service.DeleteCalendar(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "calendar deleted"})
}
