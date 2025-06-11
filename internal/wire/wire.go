//go:build wireinject
// +build wireinject

package wire

import (
	"context"

	"github.com/google/wire"
	"github.com/your-org/boilerplate-go/internal/config"
	"github.com/your-org/boilerplate-go/internal/server"
	"github.com/your-org/boilerplate-go/pkg/telemetry"
	"gorm.io/gorm"
)

// Application represents the complete application with all dependencies
type Application struct {
	Server           *server.Server
	Config           *config.Config
	CleanupTelemetry func()
}

// InitializeApplication initializes the complete application with all dependencies
func InitializeApplication(ctx context.Context) (*Application, error) {
	panic(wire.Build(
		ApplicationSet,
		wire.Struct(new(Application), "*"),
		ProvideTelemetryCleanup,
	))
}

// ProvideTelemetryCleanup provides telemetry cleanup function
func ProvideTelemetryCleanup(ctx context.Context, cfg *config.Config) func() {
	if cfg.Telemetry.Enabled {
		return telemetry.InitTelemetry(ctx, cfg)
	}
	return func() {} // Empty cleanup function when telemetry is disabled
}

// InitializeDatabase initializes only the database connection
func InitializeDatabase() (*gorm.DB, error) {
	panic(wire.Build(
		ConfigSet,
		DatabaseSet,
	))
}

// InitializeDatabaseWithConfig initializes database with provided config
func InitializeDatabaseWithConfig(cfg *config.Config) (*gorm.DB, error) {
	panic(wire.Build(
		DatabaseSet,
	))
}
