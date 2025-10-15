package db

import (
	"fmt"
	"log/slog"

	"github.com/ekosachev/go-backend-template/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect sets up a GORM Postgres connection.
func Connect(cfg *config.Config, l *slog.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode, cfg.TimeZone,
	)
	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Use a sane default GORM logger level
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	return gdb, nil
}
