package main

import (
	"context"
	"log"

	"github.com/your-org/boilerplate-go/internal/database"
	"github.com/your-org/boilerplate-go/internal/wire"
)

func main() {
	ctx := context.Background()

	// Initialize application with all dependencies using Wire
	app, err := wire.InitializeApplication(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Setup telemetry cleanup
	if app.CleanupTelemetry != nil {
		defer app.CleanupTelemetry()
	}

	// Configure database tracing if telemetry is enabled
	if app.Config.Telemetry.Enabled {
		db, err := wire.InitializeDatabaseWithConfig(app.Config)
		if err != nil {
			log.Fatalf("Failed to get database connection for tracing: %v", err)
		}
		if err := database.ConfigureTracing(db, true); err != nil {
			log.Printf("Warning: Failed to configure database tracing: %v", err)
		}
	}

	// Run migrations
	db, err := wire.InitializeDatabaseWithConfig(app.Config)
	if err != nil {
		log.Fatalf("Failed to get database connection for migrations: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Start server
	if err := app.Server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
