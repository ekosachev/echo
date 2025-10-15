package config

import (
	"errors"
	"os"
)

type Config struct {
	HTTPPort   string
	JWTSecret  string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	TimeZone   string
	Env        string
}

func Load() (*Config, error) {
	cfg := &Config{
		HTTPPort:   getEnv("HTTP_PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
		DBHost:     getEnv("DB_HOST", "db"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "appdb"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		TimeZone:   getEnv("TIME_ZONE", "UTC"),
		Env:        getEnv("ENV", "development"),
	}

	if cfg.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
