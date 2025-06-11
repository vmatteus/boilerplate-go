package config_test

import (
	"os"
	"testing"

	"github.com/your-org/boilerplate-go/internal/config"
)

func TestLoad(t *testing.T) {
	// Test loading configuration with defaults
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg == nil {
		t.Fatal("Expected config to be loaded, got nil")
	}

	// Test default values
	if cfg.Server.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", cfg.Server.Port)
	}

	if cfg.Database.Driver != "sqlite" {
		t.Errorf("Expected default driver 'sqlite', got %s", cfg.Database.Driver)
	}
}

func TestLoadWithEnvVars(t *testing.T) {
	// Set environment variable
	os.Setenv("APP_SERVER_PORT", "9090")
	defer os.Unsetenv("APP_SERVER_PORT")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Expected port from env var 9090, got %d", cfg.Server.Port)
	}
}
