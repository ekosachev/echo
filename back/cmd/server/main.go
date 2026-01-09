package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata"

	_ "github.com/ekosachev/go-backend-template/docs"
	"github.com/ekosachev/go-backend-template/internal/config"
	"github.com/ekosachev/go-backend-template/internal/db"
	"github.com/ekosachev/go-backend-template/internal/logger"
	"github.com/ekosachev/go-backend-template/internal/models"
	"github.com/ekosachev/go-backend-template/internal/router"
)

func main() {
	// Initialize logger
	l := logger.New()

	// Load config from environment
	cfg, err := config.Load()
	if err != nil {
		l.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	// Connect to DB
	gdb, err := db.Connect(cfg, l)
	if err != nil {
		l.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}

	// Apply SQL migration if `users` table does not exist
	if !gdb.Migrator().HasTable(&models.User{}) {
		l.Info("users table not found, applying migrations", slog.String("file", "migrations/0001-init.up.sql"))
		b, err := os.ReadFile("migrations/0001-init.up.sql")
		if err != nil {
			l.Error("failed to read migration file", slog.Any("error", err))
			os.Exit(1)
		}
		if res := gdb.Exec(string(b)); res.Error != nil {
			l.Error("failed to apply migration", slog.Any("error", res.Error))
			os.Exit(1)
		}
		l.Info("migration applied", slog.String("file", "migrations/0001-init.up.sql"))
	} else {
		l.Info("migration skipped, users table exists")
	}

	// Setup router and HTTP server
	r := router.NewRouter(cfg, l, gdb)

	// Graceful shutdown handling
	go func() {
		if err := r.Run(":" + cfg.HTTPPort); err != nil {
			l.Error("server run error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	l.Info("server started", slog.String("port", cfg.HTTPPort))

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	l.Info("shutting down...")

	// Close DB connection with context
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sqlDB, _ := gdb.DB()
	if err := sqlDB.Close(); err != nil {
		l.Error("failed to close db", slog.Any("error", err))
	}

	l.Info("shutdown complete")
}
