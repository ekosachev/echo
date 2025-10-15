package tests

import (
	"context"
	"testing"

	"github.com/ekosachev/go-backend-template/internal/models"
	"github.com/ekosachev/go-backend-template/internal/repository"
	"github.com/ekosachev/go-backend-template/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAuthService_RegisterAndLogin(t *testing.T) {
	// In-memory SQLite for fast tests
	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	if err := gdb.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	repo := repository.NewGormRepository[models.User](gdb)
	svc := service.NewAuthService(repo, "testsecret")

	ctx := context.Background()
	_, err = svc.Register(ctx, "Alice", "alice@example.com", "strongpassword")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Login succeeds
	token, user, err := svc.Login(ctx, "alice@example.com", "strongpassword")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if token == "" || user.Email != "alice@example.com" {
		t.Fatalf("unexpected login result")
	}

	// Login fails
	_, _, err = svc.Login(ctx, "alice@example.com", "wrong")
	if err == nil {
		t.Fatalf("expected invalid credentials error")
	}
}
