package database

import (
	"fmt"

	"github.com/your-org/boilerplate-go/internal/config"
	"github.com/your-org/boilerplate-go/internal/user/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

// Connect establishes database connection based on configuration
func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	switch cfg.Driver {
	case "postgres", "postgresql":
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Postgres.Host,
			cfg.Postgres.Port,
			cfg.Postgres.User,
			cfg.Postgres.Password,
			cfg.Postgres.Database,
			cfg.Postgres.SSLMode,
		)
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	// Add your models here for auto-migration
	// Example: return db.AutoMigrate(&models.User{}, &models.Product{})

	// Uncomment the line below to enable user migrations

	return db.AutoMigrate(&domain.User{})
}

// ConfigureTracing configures OpenTelemetry tracing for GORM
func ConfigureTracing(db *gorm.DB, enabled bool) error {
	if !enabled {
		return nil
	}

	// Add OpenTelemetry plugin to GORM
	if err := db.Use(tracing.NewPlugin()); err != nil {
		return fmt.Errorf("failed to add OpenTelemetry plugin to GORM: %w", err)
	}

	return nil
}
