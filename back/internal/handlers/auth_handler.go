package handlers

import (
	"fmt"
	"net/http"

	"github.com/ekosachev/go-backend-template/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

// Register godoc
// @Summary Register new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body registerRequest true "Register request"
// @Success 201 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}
	user, err := h.auth.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Login request"
// @Success 200 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}
	token, user, err := h.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

// Me godoc
// @Summary Get information about current user
// @Description Return current user's id email and name
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	claimsAny, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, ok := claimsAny.(map[string]any)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    claims["sub"],
		"email": claims["email"],
	})
}

type UpdateTimezoneRequest struct {
	Timezone string `json:"timezone" binding:"required"`
}

// UpdateTimezone godoc
// @Summary Update user timezone
// @Description Update current user's timezone
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateTimezoneRequest true "Update timezone request"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /api/v1/user/timezone [put]
func (h *AuthHandler) UpdateTimezone(c *gin.Context) {
	var req UpdateTimezoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	if err := h.auth.UpdateTimezone(c.Request.Context(), userID, req.Timezone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "timezone updated"})
}

// GetTimezones godoc
// @Summary Get supported timezones
// @Description Get list of all supported timezones
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Router /api/v1/timezones [get]
func (h *AuthHandler) GetTimezones(c *gin.Context) {
	timezones := h.auth.GetSupportedTimezones(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"timezones": timezones})
}

func (h *AuthHandler) getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	claimsAny, exists := c.Get("claims")
	if !exists {
		return uuid.Nil, fmt.Errorf("claims not found")
	}

	claims, ok := claimsAny.(map[string]any)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid claims format")
	}

	userIDStr, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("user ID not found in claims")
	}

	return uuid.Parse(userIDStr)
}
