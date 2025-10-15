package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekosachev/go-backend-template/internal/handlers"
	"github.com/ekosachev/go-backend-template/internal/models"
	"github.com/ekosachev/go-backend-template/internal/repository"
	"github.com/ekosachev/go-backend-template/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRouterForAuth(jwtSecret string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&models.User{})
	repo := repository.NewGormRepository[models.User](db)
	auth := service.NewAuthService(repo, jwtSecret)
	h := handlers.NewAuthHandler(auth)

	api := r.Group("/api/v1")
	api.POST("/auth/register", h.Register)
	api.POST("/auth/login", h.Login)

	return r
}

func TestRegisterAndLoginEndpoints(t *testing.T) {
	r := setupRouterForAuth("secret")

	// Register
	body := map[string]string{"name": "Bob", "email": "bob@example.com", "password": "password123"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	// Login
	b, _ = json.Marshal(map[string]string{"email": "bob@example.com", "password": "password123"})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
