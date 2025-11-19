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

func NewRouter(cfg *config.Config, l *slog.Logger, db *gorm.DB) *gin.Engine {
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

	// Wire up repositories
	usersRepo := repository.NewGormRepository[models.User](db)
	eventsRepo := repository.NewGormRepository[models.Event](db)
	calendarsRepo := repository.NewGormRepository[models.Calendar](db)

	// Wire up services
	authSvc := service.NewAuthService(usersRepo, cfg.JWTSecret)
	eventSvc := service.NewEventService(eventsRepo, usersRepo, calendarsRepo)
	calendarSvc := service.NewCalendarService(calendarsRepo)

	// Wire up handlers
	authHandler := handlers.NewAuthHandler(authSvc)
	eventHandler := handlers.NewEventHandler(eventSvc)
	calendarHandler := handlers.NewCalendarHandler(calendarSvc, eventSvc)
	telegramHandler := handlers.NewTelegramHandler(db)

	api := r.Group("/api/v1")
	{
		// Public routes
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.GET("/timezones", authHandler.GetTimezones)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Auth
			protected.GET("/auth/me", authHandler.Me)

			protected.PUT("/user/timezone", authHandler.UpdateTimezone)

			// Telegram link
			protected.POST("/telegram/link", telegramHandler.LinkTelegram)

			// Calendar routes
			protected.POST("/calendars", calendarHandler.CreateCalendar)
			protected.GET("/calendars", calendarHandler.GetUserCalendars)
			protected.GET("/calendars/:id", calendarHandler.GetCalendar)
			protected.PUT("/calendars/:id", calendarHandler.UpdateCalendar)
			protected.DELETE("/calendars/:id", calendarHandler.DeleteCalendar)

			// Event routes (with timezone support)
			protected.POST("/events", eventHandler.CreateEvent)
			protected.GET("/events/:id", eventHandler.GetEvent)
			protected.PUT("/events/:id", eventHandler.UpdateEvent)
			protected.DELETE("/events/:id", eventHandler.DeleteEvent)
		}
	}

	return r
}
