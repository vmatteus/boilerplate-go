package main

import (
	"context"
	"log"

	"github.com/your-org/boilerplate-go/internal/config"
	"github.com/your-org/boilerplate-go/internal/database"
	"github.com/your-org/boilerplate-go/internal/server"
	"github.com/your-org/boilerplate-go/pkg/logger"
	"github.com/your-org/boilerplate-go/pkg/telemetry"
)

func main() {
	ctx := context.Background()

	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	appLogger := logger.InitLogger(cfg.Logger)

	// Initialize telemetry if enabled
	var cleanupTelemetry func()
	if cfg.Telemetry.Enabled {
		cleanupTelemetry = telemetry.InitTelemetry(ctx, cfg)
		defer cleanupTelemetry()
	}

	// Initialize database
	db, err := database.Connect(cfg.Database)
	if err != nil {
		appLogger.Logger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Configure database tracing if telemetry is enabled
	if cfg.Telemetry.Enabled {
		if err := database.ConfigureTracing(db, true); err != nil {
			appLogger.Logger.Warn().Err(err).Msg("Failed to configure database tracing")
		}
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		appLogger.Logger.Fatal().Err(err).Msg("Failed to run migrations")
	}

	// Initialize and start server
	srv := server.New(cfg, db, appLogger.Logger)
	if err := srv.Start(); err != nil {
		appLogger.Logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
