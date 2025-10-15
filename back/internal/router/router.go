// router godoc
package router

import (
	"log/slog"
	"time"

	"github.com/ekosachev/go-backend-template/internal/config"
	"github.com/ekosachev/go-backend-template/internal/handlers"
	"github.com/ekosachev/go-backend-template/internal/middleware"
	"github.com/ekosachev/go-backend-template/internal/models"
	"github.com/ekosachev/go-backend-template/internal/repository"
	"github.com/ekosachev/go-backend-template/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func New(cfg *config.Config, l *slog.Logger, db *gorm.DB) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Structured logging middleware
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()

		if c.FullPath() == "/health" {
			return
		}

		l.Info("request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
			slog.Int("status", c.Writer.Status()),
			slog.String("ip", c.ClientIP()),
			slog.String("duration", time.Since(start).String()),
		)
	})

	r.Use(gin.Recovery())

	// Health endpoint
	health := handlers.NewHealthHandler(db)
	r.GET("/health", health.Health)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Wire up repositories and services
	usersRepo := repository.NewGormRepository[models.User](db)
	authSvc := service.NewAuthService(usersRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authSvc)

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			protected.GET("/me", authHandler.Me)
		}
	}

	return r
}
