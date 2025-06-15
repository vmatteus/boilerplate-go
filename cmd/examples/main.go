package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/your-org/boilerplate-go/internal/config"
	"github.com/your-org/boilerplate-go/internal/logger"
	"github.com/your-org/boilerplate-go/internal/user/examples"
)

func main() {
	fmt.Println("🚀 Advanced Logger with OpenTelemetry - Demo")
	fmt.Println(strings.Repeat("=", 60))

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Initialize logger
	appLogger := logger.InitLogger(cfg.Logger)

	ctx := context.Background()

	// Create logger examples instance
	loggerExamples := examples.NewLoggerExamples(&appLogger)

	// Demonstrate basic logging
	fmt.Println("\n📝 1. Basic Logging:")
	loggerExamples.ExampleBasicLogging(ctx)

	// Demonstrate logging with tracing
	fmt.Println("\n🔍 2. Logging with OpenTelemetry Tracing:")
	loggerExamples.ExampleWithTracing(ctx, 12345)

	// Demonstrate error handling
	fmt.Println("\n❌ 3. Error Handling:")
	if err := loggerExamples.ExampleErrorHandling(ctx); err != nil {
		appLogger.LogInfo(ctx, "Error handling example completed as expected", map[string]interface{}{
			"error_expected": true,
		})
	}

	// Demonstrate structured logging
	fmt.Println("\n📊 4. Structured Logging:")
	loggerExamples.ExampleStructuredLogging(ctx)

	// Demonstrate contextual logging
	fmt.Println("\n🔗 5. Contextual Logging:")
	loggerExamples.ExampleContextualLogging(ctx, "req_abc123", "user_12345")

	fmt.Println("\n✅ Demo completed! Check the logs above.")
	fmt.Println("\nTry different configurations:")
	fmt.Println("- make run-debug    (debug level)")
	fmt.Println("- make run-json     (JSON format)")
	fmt.Println("- make run-file     (file output)")
}
