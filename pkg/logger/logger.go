package logger

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/your-org/boilerplate-go/internal/config"
)

// Logger wraps zerolog.Logger with additional functionality
type Logger struct {
	Logger    zerolog.Logger
	AppLogger AppLogger
}

// InitLogger initializes the logger based on configuration
func InitLogger(cfg config.LoggerConfig) Logger {
	// Set log level
	level := parseLogLevel(cfg.Level)
	zerolog.SetGlobalLevel(level)

	var logger zerolog.Logger
	var appLogger AppLogger

	// Initialize the appropriate logger based on provider
	switch cfg.Provider {
	case "stdout", "":
		appLogger = NewStdoutLogger(cfg)
		logger = initZerologLogger(cfg)
	case "file":
		appLogger = NewFileLogger(cfg)
		logger = initZerologLogger(cfg)
	case "elasticsearch":
		appLogger = NewElasticsearchLogger(cfg)
		logger = initZerologLogger(cfg)
	default:
		fmt.Printf("Invalid logger provider '%s', falling back to stdout\n", cfg.Provider)
		appLogger = NewStdoutLogger(cfg)
		logger = initZerologLogger(cfg)
	}

	// Set as global logger
	log.Logger = logger

	return Logger{
		Logger:    logger,
		AppLogger: appLogger,
	}
}

// initZerologLogger initializes the zerolog logger for console output
func initZerologLogger(cfg config.LoggerConfig) zerolog.Logger {
	var logger zerolog.Logger

	if cfg.Format == "json" {
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	} else {
		// Console format with colors
		output := zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			FormatLevel: func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
			},
			FormatMessage: func(i interface{}) string {
				return fmt.Sprintf("***%s****", i)
			},
			FormatFieldName: func(i interface{}) string {
				return fmt.Sprintf("%s:", i)
			},
			FormatFieldValue: func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("%s", i))
			},
		}
		logger = zerolog.New(output).With().Timestamp().Logger()
	}

	return logger
}

// Log logs a message with the specified level and fields
func (l *Logger) Log(ctx context.Context, level, message string, fields map[string]interface{}) {
	// Log to zerolog (console/stderr)
	switch level {
	case "debug":
		l.Logger.Debug().Fields(fields).Msg(message)
	case "info":
		l.Logger.Info().Fields(fields).Msg(message)
	case "warn":
		l.Logger.Warn().Fields(fields).Msg(message)
	case "error":
		l.Logger.Error().Fields(fields).Msg(message)
	case "fatal":
		l.Logger.Fatal().Fields(fields).Msg(message)
	default:
		l.Logger.Info().Fields(fields).Msg(message)
	}

	// Log to configured provider (file/elasticsearch)
	if l.AppLogger != nil {
		l.AppLogger.Log(ctx, level, message, fields)
	}
}

// parseLogLevel converts string log level to zerolog level
func parseLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

// WithTraceID adds trace ID to logger context
func WithTraceID(logger zerolog.Logger, traceID string) zerolog.Logger {
	return logger.With().Str("trace_id", traceID).Logger()
}

// WithSpanID adds span ID to logger context
func WithSpanID(logger zerolog.Logger, spanID string) zerolog.Logger {
	return logger.With().Str("span_id", spanID).Logger()
}

// WithRequestID adds request ID to logger context
func WithRequestID(logger zerolog.Logger, requestID string) zerolog.Logger {
	return logger.With().Str("request_id", requestID).Logger()
}

// WithUserID adds user ID to logger context
func WithUserID(logger zerolog.Logger, userID string) zerolog.Logger {
	return logger.With().Str("user_id", userID).Logger()
}
