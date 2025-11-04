package handlers

import (
	"net/http"

	"github.com/ekosachev/go-backend-template/internal/models"
	"github.com/ekosachev/go-backend-template/internal/service"
	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) CreateEvent(ctx *gin.Context) {
	var event models.Event

	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	if err := h.service.CreateEvent(ctx.Request.Context(), &event); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, event)
}

func (h *EventHandler) GetEventByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "event ID is empty"})
		return
	}

	event, err := h.service.GetEventByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, event)
}

func (h *EventHandler) UpdateEvent(ctx *gin.Context) {
	var event models.Event

	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	if err := h.service.UpdateEvent(ctx.Request.Context(), &event); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, event)
}

func (h *EventHandler) DeleteEvent(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "event ID is empty"})
		return
	}

	if err := h.service.DeleteEvent(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "event deleted"})
}
