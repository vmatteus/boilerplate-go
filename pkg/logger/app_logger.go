package logger

import "context"

// AppLogger interface defines the contract for different logger implementations
type AppLogger interface {
	Log(ctx context.Context, level, message string, fields map[string]interface{})
	AddField(key string, value interface{})
}
