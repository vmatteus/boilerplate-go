package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/your-org/boilerplate-go/internal/config"
)

type FileLogger struct {
	config config.LoggerConfig
	fields map[string]interface{}
}

func NewFileLogger(cfg config.LoggerConfig) *FileLogger {
	return &FileLogger{
		config: cfg,
		fields: make(map[string]interface{}),
	}
}

func (l *FileLogger) AddField(key string, value interface{}) {
	l.fields[key] = value
}

func (l *FileLogger) Log(ctx context.Context, level, message string, fields map[string]interface{}) {
	// Ensure log directory exists
	logDir := filepath.Dir(l.config.Filepath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
		return
	}

	// Open log file for writing
	file, err := os.OpenFile(l.config.Filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer file.Close()

	// Create logger with JSON format for file
	logger := zerolog.New(file).With().Timestamp().Logger()

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
		event = logger.Debug()
	case "info":
		event = logger.Info()
	case "warn":
		event = logger.Warn()
	case "error":
		event = logger.Error()
	case "fatal":
		event = logger.Fatal()
	default:
		event = logger.Info()
	}

	// Add fields to event
	for key, value := range mergedFields {
		event = event.Interface(key, value)
	}

	// Log the message
	event.Msg(message)
}
