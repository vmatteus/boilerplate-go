package examples

import (
	"context"

	"github.com/your-org/boilerplate-go/internal/logger"
)

// RunAllExamples executes all available examples
func RunAllExamples() {
	ctx := context.Background()
	appLogger := getLogger()

	appLogger.LogInfo(ctx, "Starting All Examples", map[string]interface{}{
		"total_categories": 2,
		"status":           "starting",
	})

	// Run EventBus Examples
	runEventBusExamples(ctx, appLogger)

	// Run Logger Examples
	runLoggerExamples(ctx, appLogger)

	appLogger.LogInfo(ctx, "All examples executed successfully", map[string]interface{}{
		"status": "completed",
	})
}

// RunEventBusExamples executes only EventBus-related examples
func RunEventBusExamples() {
	ctx := context.Background()
	appLogger := getLogger()

	runEventBusExamples(ctx, appLogger)
}

// RunLoggerExamples executes only Logger-related examples
func RunLoggerExamples() {
	ctx := context.Background()
	appLogger := getLogger()

	runLoggerExamples(ctx, appLogger)
}

func runEventBusExamples(ctx context.Context, appLogger *logger.Logger) {
	appLogger.LogInfo(ctx, "Running EventBus Examples", map[string]interface{}{
		"category":       "eventbus",
		"total_examples": 6,
		"status":         "starting",
	})

	NewEventBusExample()
	AsyncPublishingExample()
	ContextCancellationExample()
	BufferOverflowExample()
	SubscribeOnceExample()
	ChannelEventBusExample()

	appLogger.LogInfo(ctx, "EventBus examples completed", map[string]interface{}{
		"category": "eventbus",
		"status":   "completed",
	})
}

func runLoggerExamples(ctx context.Context, appLogger *logger.Logger) {
	appLogger.LogInfo(ctx, "Running Logger Examples", map[string]interface{}{
		"category":       "logger",
		"total_examples": 5,
		"status":         "starting",
	})

	// Create logger examples instance
	loggerExamples := NewLoggerExamples(appLogger)

	// Demonstrate basic logging
	appLogger.LogInfo(ctx, "Example: Basic Logging", map[string]interface{}{
		"example": "basic_logging",
	})
	loggerExamples.ExampleBasicLogging(ctx)

	// Demonstrate logging with tracing
	appLogger.LogInfo(ctx, "Example: Logging with OpenTelemetry Tracing", map[string]interface{}{
		"example": "tracing_logging",
	})
	loggerExamples.ExampleWithTracing(ctx, 12345)

	// Demonstrate error handling
	appLogger.LogInfo(ctx, "Example: Error Handling", map[string]interface{}{
		"example": "error_handling",
	})
	if err := loggerExamples.ExampleErrorHandling(ctx); err != nil {
		appLogger.LogInfo(ctx, "Error handling example completed as expected", map[string]interface{}{
			"error_expected": true,
		})
	}

	// Demonstrate structured logging
	appLogger.LogInfo(ctx, "Example: Structured Logging", map[string]interface{}{
		"example": "structured_logging",
	})
	loggerExamples.ExampleStructuredLogging(ctx)

	// Demonstrate contextual logging
	appLogger.LogInfo(ctx, "Example: Contextual Logging", map[string]interface{}{
		"example": "contextual_logging",
	})
	loggerExamples.ExampleContextualLogging(ctx, "req_abc123", "user_12345")

	appLogger.LogInfo(ctx, "Logger examples completed", map[string]interface{}{
		"category": "logger",
		"status":   "completed",
	})
}

func init() {
	// Initialize any global configurations if needed
}
