package examples

import (
	"context"
	"errors"
	"time"

	"github.com/your-org/boilerplate-go/internal/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// LoggerExamples demonstrates best practices for using the improved logger
type LoggerExamples struct {
	logger *logger.Logger
}

// NewLoggerExamples creates a new LoggerExamples instance
func NewLoggerExamples(logger *logger.Logger) *LoggerExamples {
	return &LoggerExamples{
		logger: logger,
	}
}

// ExampleBasicLogging shows basic logging with different levels
func (e *LoggerExamples) ExampleBasicLogging(ctx context.Context) {
	// Basic info logging
	e.logger.LogInfo(ctx, "Application started successfully")

	// Logging with fields
	e.logger.LogInfo(ctx, "User operation completed", map[string]interface{}{
		"user_id":    12345,
		"operation":  "create",
		"duration":   "150ms",
		"ip_address": "192.168.1.100",
	})

	// Warning with context
	e.logger.LogWarn(ctx, "Rate limit approaching", map[string]interface{}{
		"current_requests": 950,
		"limit":            1000,
		"window":           "1m",
	})

	// Error logging with error object
	err := errors.New("database connection failed")
	e.logger.LogError(ctx, "Failed to connect to database", err, map[string]interface{}{
		"database_host": "localhost:5432",
		"retry_count":   3,
	})

	// Debug logging (only shown when debug level is enabled)
	e.logger.LogDebug(ctx, "Cache hit", map[string]interface{}{
		"cache_key": "user:12345",
		"ttl":       "300s",
	})
}

// ExampleWithTracing shows logging combined with OpenTelemetry tracing
func (e *LoggerExamples) ExampleWithTracing(ctx context.Context, userID int64) {
	// Start a new span for this operation
	ctx, span := otel.Tracer("user-examples").Start(ctx, "ExampleWithTracing")
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(
		attribute.Int64("user.id", userID),
		attribute.String("operation", "example_tracing"),
	)

	start := time.Now()

	e.logger.LogInfo(ctx, "Starting traced operation", map[string]interface{}{
		"user_id":   userID,
		"operation": "example_tracing",
	})

	// Simulate some work with child spans
	e.simulateWork(ctx, userID)

	duration := time.Since(start)

	// Add duration to span
	span.SetAttributes(attribute.Int64("duration_ms", duration.Milliseconds()))

	e.logger.LogInfo(ctx, "Traced operation completed", map[string]interface{}{
		"user_id":  userID,
		"duration": duration.Milliseconds(),
	})
}

// simulateWork creates child spans and demonstrates nested tracing
func (e *LoggerExamples) simulateWork(ctx context.Context, userID int64) {
	// Database operation simulation
	ctx, dbSpan := otel.Tracer("user-examples").Start(ctx, "database.query")
	defer dbSpan.End()

	dbSpan.SetAttributes(
		attribute.String("db.operation", "SELECT"),
		attribute.String("db.table", "users"),
		attribute.Int64("user.id", userID),
	)

	e.logger.LogDebug(ctx, "Executing database query", map[string]interface{}{
		"query": "SELECT * FROM users WHERE id = ?",
		"args":  []interface{}{userID},
	})

	// Simulate database work
	time.Sleep(50 * time.Millisecond)

	dbSpan.SetStatus(codes.Ok, "Query executed successfully")

	// Cache operation simulation
	ctx, cacheSpan := otel.Tracer("user-examples").Start(ctx, "cache.set")
	defer cacheSpan.End()

	cacheSpan.SetAttributes(
		attribute.String("cache.key", "user:profile"),
		attribute.Int64("user.id", userID),
	)

	e.logger.LogDebug(ctx, "Caching user profile", map[string]interface{}{
		"cache_key": "user:profile",
		"user_id":   userID,
		"ttl":       "300s",
	})

	time.Sleep(10 * time.Millisecond)
	cacheSpan.SetStatus(codes.Ok, "Cache updated")
}

// ExampleErrorHandling demonstrates error logging with tracing
func (e *LoggerExamples) ExampleErrorHandling(ctx context.Context) error {
	ctx, span := otel.Tracer("user-examples").Start(ctx, "ExampleErrorHandling")
	defer span.End()

	e.logger.LogInfo(ctx, "Starting error handling example")

	// Simulate an operation that might fail
	if err := e.simulateFailingOperation(ctx); err != nil {
		// Record the error in the span
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		// Log the error with context
		e.logger.LogError(ctx, "Operation failed", err, map[string]interface{}{
			"operation": "simulate_failing_operation",
			"retry":     false,
		})

		return err
	}

	e.logger.LogInfo(ctx, "Error handling example completed successfully")
	return nil
}

// simulateFailingOperation simulates an operation that might fail
func (e *LoggerExamples) simulateFailingOperation(ctx context.Context) error {
	ctx, span := otel.Tracer("user-examples").Start(ctx, "simulateFailingOperation")
	defer span.End()

	e.logger.LogDebug(ctx, "Simulating potentially failing operation")

	// Simulate some work
	time.Sleep(25 * time.Millisecond)

	// Simulate a failure
	err := errors.New("simulated operation failure")

	e.logger.LogWarn(ctx, "Operation encountered an error", map[string]interface{}{
		"error_type":  "simulation",
		"recoverable": false,
	})

	return err
}

// ExampleStructuredLogging shows advanced structured logging patterns
func (e *LoggerExamples) ExampleStructuredLogging(ctx context.Context) {
	// Business event logging
	e.logger.LogInfo(ctx, "User registration completed", map[string]interface{}{
		"event_type":          "user_registration",
		"user_id":             12345,
		"email":               "user@example.com",
		"registration_source": "web",
		"country":             "BR",
		"timestamp":           time.Now().Unix(),
	})

	// Performance metrics logging
	e.logger.LogInfo(ctx, "API endpoint performance", map[string]interface{}{
		"endpoint":      "/api/v1/users",
		"method":        "POST",
		"response_time": 156,
		"status_code":   201,
		"request_size":  1024,
		"response_size": 512,
	})

	// Security event logging
	e.logger.LogWarn(ctx, "Suspicious activity detected", map[string]interface{}{
		"event_type":      "security_alert",
		"ip_address":      "192.168.1.100",
		"user_agent":      "Mozilla/5.0...",
		"failed_attempts": 5,
		"time_window":     "5m",
		"action":          "account_locked",
	})

	// Integration logging
	e.logger.LogInfo(ctx, "External API call", map[string]interface{}{
		"service":       "payment_gateway",
		"endpoint":      "/api/v2/charge",
		"request_id":    "req_abc123",
		"response_time": 1250,
		"status":        "success",
		"amount":        99.99,
		"currency":      "BRL",
	})
}

// ExampleContextualLogging demonstrates logging with rich context
func (e *LoggerExamples) ExampleContextualLogging(ctx context.Context, requestID, userID string) {
	// Create a context-aware logger session
	// The logger will automatically extract trace information from context

	e.logger.LogInfo(ctx, "Processing user request", map[string]interface{}{
		"request_id": requestID,
		"user_id":    userID,
		"handler":    "user_controller",
	})

	// Simulate processing steps with context
	for i, step := range []string{"validate", "authorize", "process", "respond"} {
		ctx, stepSpan := otel.Tracer("user-examples").Start(ctx, "step."+step)

		stepSpan.SetAttributes(
			attribute.String("step.name", step),
			attribute.Int("step.index", i+1),
		)

		e.logger.LogDebug(ctx, "Processing step", map[string]interface{}{
			"step":       step,
			"step_index": i + 1,
			"request_id": requestID,
		})

		// Simulate step processing time
		time.Sleep(time.Duration(10+i*5) * time.Millisecond)

		stepSpan.End()
	}

	e.logger.LogInfo(ctx, "Request processing completed", map[string]interface{}{
		"request_id": requestID,
		"user_id":    userID,
		"status":     "success",
	})
}
