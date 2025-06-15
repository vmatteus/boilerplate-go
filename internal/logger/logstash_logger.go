package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/your-org/boilerplate-go/internal/config"
)

type LogstashLogger struct {
	config     config.LoggerConfig
	connection net.Conn
	fields     map[string]interface{}
}

func NewLogstashLogger(cfg config.LoggerConfig) *LogstashLogger {
	logger := &LogstashLogger{
		config: cfg,
		fields: make(map[string]interface{}),
	}

	// Try to establish connection to Logstash
	if cfg.Url != "" {
		conn, err := net.Dial("tcp", cfg.Url)
		if err != nil {
			fmt.Printf("Failed to connect to Logstash at %s: %v\n", cfg.Url, err)
		} else {
			logger.connection = conn
		}
	}

	return logger
}

func (l *LogstashLogger) AddField(key string, value interface{}) {
	l.fields[key] = value
}

func (l *LogstashLogger) Log(ctx context.Context, level, message string, fields map[string]interface{}) {
	// If no connection, skip logging
	if l.connection == nil {
		return
	}

	// Merge stored fields with provided fields
	mergedFields := make(map[string]interface{})
	for k, v := range l.fields {
		mergedFields[k] = v
	}
	for k, v := range fields {
		mergedFields[k] = v
	}

	// Create log entry
	logEntry := map[string]interface{}{
		"@timestamp": time.Now().UTC().Format(time.RFC3339),
		"level":      level,
		"message":    message,
		"fields":     mergedFields,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		fmt.Printf("Failed to marshal log entry to JSON: %v\n", err)
		return
	}

	// Send to Logstash
	_, err = l.connection.Write(append(jsonData, '\n'))
	if err != nil {
		fmt.Printf("Failed to send log to Logstash: %v\n", err)
	}
}

func (l *LogstashLogger) Close() error {
	if l.connection != nil {
		return l.connection.Close()
	}
	return nil
}
