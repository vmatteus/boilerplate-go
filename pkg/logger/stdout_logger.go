package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/your-org/boilerplate-go/internal/config"
)

type StdoutLogger struct {
	config config.LoggerConfig
	logger zerolog.Logger
	fields map[string]interface{}
}

func NewStdoutLogger(cfg config.LoggerConfig) *StdoutLogger {
	var logger zerolog.Logger

	if cfg.Format == "json" {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		// Console format with colors
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05",
		}
		logger = zerolog.New(output).With().Timestamp().Logger()
	}

	return &StdoutLogger{
		config: cfg,
		logger: logger,
		fields: make(map[string]interface{}),
	}
}

func (l *StdoutLogger) AddField(key string, value interface{}) {
	l.fields[key] = value
}

func (l *StdoutLogger) Log(ctx context.Context, level, message string, fields map[string]interface{}) {
	// Merge stored fields with provided fields
	mergedFields := make(map[string]interface{})
	for k, v := range l.fields {
		mergedFields[k] = v
	}
	for k, v := range fields {
		mergedFields[k] = v
	}

	// Create log event based on level
	var event *zerolog.Event
	switch level {
	case "debug":
		event = l.logger.Debug()
	case "info":
		event = l.logger.Info()
	case "warn":
		event = l.logger.Warn()
	case "error":
		event = l.logger.Error()
	case "fatal":
		event = l.logger.Fatal()
	default:
		event = l.logger.Info()
	}

	// Add fields to event
	for key, value := range mergedFields {
		event = event.Interface(key, value)
	}

	// Log the message
	event.Msg(message)
}
