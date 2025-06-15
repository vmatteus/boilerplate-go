package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/your-org/boilerplate-go/internal/config"
	"github.com/your-org/boilerplate-go/internal/logger"
	"github.com/your-org/boilerplate-go/internal/middleware"
	"github.com/your-org/boilerplate-go/internal/user/presentation"
	"gorm.io/gorm"
)

type Server struct {
	config         *config.Config
	db             *gorm.DB
	logger         *logger.Logger
	router         *gin.Engine
	userController *presentation.UserController
}

// New creates a new server instance
func New(cfg *config.Config, db *gorm.DB, appLogger *logger.Logger, userController *presentation.UserController) *Server {
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create router
	router := gin.New()

	return &Server{
		config:         cfg,
		db:             db,
		logger:         appLogger,
		router:         router,
		userController: userController,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Setup middleware
	s.setupMiddleware()

	// Setup routes
	s.setupRoutes()

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port),
		Handler: s.router,
	}

	// Start server in a goroutine
	go func() {
		ctx := context.Background()
		s.logger.LogInfo(ctx, "Starting HTTP server", map[string]interface{}{
			"host": s.config.Server.Host,
			"port": s.config.Server.Port,
		})

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.LogError(ctx, "Failed to start server", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx := context.Background()
	s.logger.LogInfo(ctx, "Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		s.logger.LogError(ctx, "Server forced to shutdown", err)
	}

	s.logger.LogInfo(ctx, "Server exited")
	return nil
}

// setupMiddleware configures middleware
func (s *Server) setupMiddleware() {
	s.router.Use(middleware.Logger(s.logger.Logger))
	s.router.Use(middleware.Recovery(s.logger.Logger))
	s.router.Use(middleware.CORS())

	// Add OpenTelemetry middleware if telemetry is enabled
	if s.config.Telemetry.Enabled {
		s.router.Use(middleware.OpenTelemetry(s.config.Application.Name))
	}
}

// setupRoutes configures API routes
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.GET("/health", s.healthCheck)

	// API routes
	v1 := s.router.Group("/api/v1")
	{
		// Welcome endpoint
		v1.GET("/", s.welcome)

		// User routes - using injected controller
		s.userController.RegisterRoutes(v1)
	}
}

// healthCheck handles health check requests
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Server is running",
	})
}

// welcome handles welcome requests
func (s *Server) welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to Boilerplate Go API",
		"version": "1.0.0",
	})
}
