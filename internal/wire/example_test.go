package wire_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/your-org/boilerplate-go/internal/wire"
)

// TestWireConfigInitialization tests only the configuration initialization
func TestWireConfigInitialization(t *testing.T) {
	// Test only the config provider which doesn't require database
	cfg, err := wire.ProvideConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Verify that configuration was loaded correctly
	assert.Equal(t, "boilerplate-go", cfg.Application.Name)
	assert.Equal(t, "1.0.0", cfg.Application.Version)
}

// TestWireProviders tests individual providers without full application initialization
func TestWireProviders(t *testing.T) {
	// Test config provider
	cfg, err := wire.ProvideConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Test logger provider
	logger := wire.ProvideLogger(cfg)
	assert.NotNil(t, logger.Logger)

	// Test zerolog provider
	zerologLogger := wire.ProvideZerologLogger(logger)
	assert.NotNil(t, zerologLogger)
}

// ExampleInitializeApplication shows how to use Wire in your main function
func ExampleInitializeApplication() {
	ctx := context.Background()

	// Initialize the complete application with all dependencies
	app, err := wire.InitializeApplication(ctx)
	if err != nil {
		panic(err)
	}

	// Setup cleanup
	defer func() {
		if app.CleanupTelemetry != nil {
			app.CleanupTelemetry()
		}
	}()

	// Start your server
	// app.Server.Start()
}
