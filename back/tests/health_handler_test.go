package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekosachev/go-backend-template/internal/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	h := handlers.NewHealthHandler(db)
	r.GET("/health", h.Health)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
