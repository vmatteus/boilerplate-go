package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// Logger middleware logs HTTP requests with OpenTelemetry integration
func Logger(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Generate request ID if not present
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Header("X-Request-ID", requestID)
		}

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get trace information if available
		spanCtx := trace.SpanContextFromContext(c.Request.Context())

		// Build log event
		logEvent := logger.Info().
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", c.Writer.Status()).
			Dur("latency", latency).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Str("request_id", requestID).
			Int("body_size", c.Writer.Size())

		// Add trace information if available
		if spanCtx.IsValid() {
			logEvent.Str("trace_id", spanCtx.TraceID().String()).
				Str("span_id", spanCtx.SpanID().String())
		}

		// Add query parameters if present
		if raw != "" {
			logEvent.Str("query", raw)
		}

		// Add error information if status >= 400
		if c.Writer.Status() >= 400 {
			if len(c.Errors) > 0 {
				logEvent.Str("error", c.Errors.String())
			}
		}

		logEvent.Msg("HTTP Request")
	}
}

// Recovery middleware recovers from panics with enhanced logging
func Recovery(logger zerolog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Get trace information if available
		spanCtx := trace.SpanContextFromContext(c.Request.Context())

		logEvent := logger.Error().
			Interface("error", recovered).
			Str("path", c.Request.URL.Path).
			Str("method", c.Request.Method).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent())

		// Add trace information if available
		if spanCtx.IsValid() {
			logEvent.Str("trace_id", spanCtx.TraceID().String()).
				Str("span_id", spanCtx.SpanID().String())
		}

		logEvent.Msg("Panic recovered")

		c.AbortWithStatus(500)
	})
}
