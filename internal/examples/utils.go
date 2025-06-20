package examples

import (
	"github.com/your-org/boilerplate-go/internal/config"
	"github.com/your-org/boilerplate-go/internal/logger"
)

// getLogger creates and returns a configured logger instance
// This function is shared across all example files
func getLogger() *logger.Logger {
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}
	appLogger := logger.InitLogger(cfg.Logger)
	return &appLogger
}

// Common configuration for examples
const (
	DefaultChannelBuffer   = 100
	DefaultProcessingDelay = 500 // milliseconds
	DefaultTimeoutDuration = 5   // seconds
)
