package logger

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/your-org/boilerplate-go/internal/config"
	"go.opentelemetry.io/otel/trace"
)

// Logger wraps zerolog.Logger with additional functionality
type Logger struct {
	Logger    zerolog.Logger
	AppLogger AppLogger
}

// InitLogger initializes the logger based on configuration
func InitLogger(cfg config.LoggerConfig) Logger {
	// Set zerolog error stack marshaler for better error handling
	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		return err.Error()
	}

	var logger zerolog.Logger
	var appLogger AppLogger

	switch cfg.Provider {
	case "stdout":
		appLogger = NewStdoutLogger(cfg)
	case "file":
		appLogger = NewFileLogger(cfg)
	case "logstash":
		appLogger = NewLogstashLogger(cfg)
	case "elasticsearch":
		appLogger = NewElasticsearchLogger(cfg)
	default:
		log.Fatal().Msg("Invalid logger provider specified")
	}

	return Logger{
		Logger:    logger,
		AppLogger: appLogger,
	}
}

// Log logs a message with the specified level and fields
func (l *Logger) Log(ctx context.Context, level, message string, fields map[string]interface{}) {
	// Enrich fields with OpenTelemetry trace information
	enrichedFields := l.enrichWithTraceInfo(ctx, fields)

	switch level {
	case "info":
		l.Logger.Info().Fields(enrichedFields).Msg(message)
	case "warn":
		l.Logger.Warn().Fields(enrichedFields).Msg(message)
	case "error":
		l.Logger.Error().Fields(enrichedFields).Msg(message)
	case "debug":
		l.Logger.Debug().Fields(enrichedFields).Msg(message)
	case "fatal":
		l.Logger.Fatal().Fields(enrichedFields).Msg(message)
	default:
		l.Logger.Debug().Fields(enrichedFields).Msg(message)
	}

	l.AppLogger.Log(ctx, level, message, enrichedFields)
}

// enrichWithTraceInfo adds OpenTelemetry trace information to log fields
func (l *Logger) enrichWithTraceInfo(ctx context.Context, fields map[string]interface{}) map[string]interface{} {
	enriched := make(map[string]interface{})

	// Copy existing fields
	for k, v := range fields {
		enriched[k] = v
	}

	// Add OpenTelemetry trace information
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		enriched["trace_id"] = span.SpanContext().TraceID().String()
		enriched["span_id"] = span.SpanContext().SpanID().String()

		if span.SpanContext().IsSampled() {
			enriched["sampled"] = true
		}
	}

	return enriched
}

// LogInfo logs an info message with optional fields
func (l *Logger) LogInfo(ctx context.Context, message string, fields ...map[string]interface{}) {
	allFields := make(map[string]interface{})
	for _, fieldMap := range fields {
		for k, v := range fieldMap {
			allFields[k] = v
		}
	}
	l.Log(ctx, "info", message, allFields)
}

// LogError logs an error message with optional fields
func (l *Logger) LogError(ctx context.Context, message string, err error, fields ...map[string]interface{}) {
	allFields := make(map[string]interface{})
	for _, fieldMap := range fields {
		for k, v := range fieldMap {
			allFields[k] = v
		}
	}
	if err != nil {
		allFields["error"] = err.Error()
	}
	l.Log(ctx, "error", message, allFields)
}

// LogWarn logs a warning message with optional fields
func (l *Logger) LogWarn(ctx context.Context, message string, fields ...map[string]interface{}) {
	allFields := make(map[string]interface{})
	for _, fieldMap := range fields {
		for k, v := range fieldMap {
			allFields[k] = v
		}
	}
	l.Log(ctx, "warn", message, allFields)
}

// LogDebug logs a debug message with optional fields
func (l *Logger) LogDebug(ctx context.Context, message string, fields ...map[string]interface{}) {
	allFields := make(map[string]interface{})
	for _, fieldMap := range fields {
		for k, v := range fieldMap {
			allFields[k] = v
		}
	}
	l.Log(ctx, "debug", message, allFields)
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
