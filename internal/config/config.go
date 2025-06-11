package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Logger      LoggerConfig      `mapstructure:"logger"`
	Telemetry   TelemetryConfig   `mapstructure:"telemetry"`
	Application ApplicationConfig `mapstructure:"application"`
	Apm         Apm               `mapstructure:"apm"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release, test
}

type DatabaseConfig struct {
	Driver   string           `mapstructure:"driver"`
	Postgres PostgreSQLConfig `mapstructure:"postgres"`
	SQLite   SQLiteConfig     `mapstructure:"sqlite"`
}

type PostgreSQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"sslmode"`
}

type SQLiteConfig struct {
	Path string `mapstructure:"path"`
}

type LoggerConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`   // json, console
	Provider string `mapstructure:"provider"` // stdout, file, elasticsearch
	Index    string `mapstructure:"index"`    // for elasticsearch
	Url      string `mapstructure:"url"`      // for elasticsearch
	ApiKey   string `mapstructure:"api_key"`  // for elasticsearch
	Username string `mapstructure:"username"` // for elasticsearch
	Password string `mapstructure:"password"` // for elasticsearch
	Filepath string `mapstructure:"filepath"` // for file logging
}

type TelemetryConfig struct {
	Enabled               bool   `mapstructure:"enabled"`
	TracingEnabled        bool   `mapstructure:"tracing_enabled"`
	MetricsEnabled        bool   `mapstructure:"metrics_enabled"`
	HostMetricsEnabled    bool   `mapstructure:"host_metrics_enabled"`
	RuntimeMetricsEnabled bool   `mapstructure:"runtime_metrics_enabled"`
	Endpoint              string `mapstructure:"endpoint"`
	Headers               string `mapstructure:"headers"`
	Attributes            string `mapstructure:"attributes"`
}

type ApplicationConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

type Apm struct {
	Url            string `mapstructure:"url"`
	ServiceName    string `mapstructure:"service_name"`
	ServiceVersion string `mapstructure:"service_version"`
	Attributes     string `mapstructure:"attributes"`
	Environment    string `mapstructure:"environment"`
	Headers        string `mapstructure:"headers"`
}

// Load reads configuration from file and environment variables
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("$HOME/.boilerplate-go")

	// Set default values
	setDefaults()

	// Environment variables
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")

	// Database defaults
	viper.SetDefault("database.driver", "sqlite")
	viper.SetDefault("database.postgres.host", "localhost")
	viper.SetDefault("database.postgres.port", 5432)
	viper.SetDefault("database.postgres.user", "postgres")
	viper.SetDefault("database.postgres.password", "password")
	viper.SetDefault("database.postgres.database", "boilerplate")
	viper.SetDefault("database.postgres.sslmode", "disable")
	viper.SetDefault("database.sqlite.path", "./data/app.db")

	// Logger defaults
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.format", "console")
	viper.SetDefault("logger.provider", "stdout")
	viper.SetDefault("logger.filepath", "./logs/app.log")
	viper.SetDefault("logger.index", "boilerplate-go-logs")
	viper.SetDefault("logger.url", "http://localhost:9200")
	viper.SetDefault("logger.username", "")
	viper.SetDefault("logger.password", "")
	viper.SetDefault("logger.api_key", "")

	// Telemetry defaults
	viper.SetDefault("telemetry.enabled", false)
	viper.SetDefault("telemetry.tracing_enabled", true)
	viper.SetDefault("telemetry.metrics_enabled", true)
	viper.SetDefault("telemetry.host_metrics_enabled", false)
	viper.SetDefault("telemetry.runtime_metrics_enabled", false)
	viper.SetDefault("telemetry.endpoint", "")
	viper.SetDefault("telemetry.headers", "")
	viper.SetDefault("telemetry.attributes", "")

	// Application defaults
	viper.SetDefault("application.name", "boilerplate-go")
	viper.SetDefault("application.version", "1.0.0")
	viper.SetDefault("application.environment", "development")
}
