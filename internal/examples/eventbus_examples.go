package examples

import (
	"context"
	"time"

	"github.com/your-org/boilerplate-go/pkg/events"
)

type UserCreatedEvent struct {
	events.BaseEvent
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// NewEventBusExample demonstrates the EventBus with multiple subscribers
func NewEventBusExample() {
	ctx := context.Background()
	appLogger := getLogger()

	appLogger.LogInfo(ctx, "EventBus Channel-Based Example", map[string]interface{}{
		"example_type": "eventbus_channels",
		"status":       "starting",
	})

	eventBus := events.NewEventBus()

	// Email Service Subscriber
	emailChan := make(chan events.Event, DefaultChannelBuffer)
	eventBus.Subscribe("user.created", emailChan)

	go func() {
		for event := range emailChan {
			user := event.(*UserCreatedEvent)
			appLogger.LogInfo(ctx, "Email service processing", map[string]interface{}{
				"service":  "email",
				"action":   "send_welcome_email",
				"username": user.Username,
				"email":    user.Email,
			})
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Analytics Service Subscriber
	analyticsChan := make(chan events.Event, DefaultChannelBuffer)
	eventBus.Subscribe("user.created", analyticsChan)

	go func() {
		for event := range analyticsChan {
			user := event.(*UserCreatedEvent)
			appLogger.LogInfo(ctx, "Analytics service processing", map[string]interface{}{
				"service": "analytics",
				"action":  "record_user_metrics",
				"user_id": user.UserID,
			})
			time.Sleep(200 * time.Millisecond)
		}
	}()

	// Notification Service Subscriber
	notificationChan := make(chan events.Event, DefaultChannelBuffer)
	eventBus.Subscribe("user.created", notificationChan)

	go func() {
		for event := range notificationChan {
			user := event.(*UserCreatedEvent)
			appLogger.LogInfo(ctx, "Notification service processing", map[string]interface{}{
				"service":  "notifications",
				"action":   "create_welcome_notification",
				"username": user.Username,
			})
			time.Sleep(300 * time.Millisecond)
		}
	}()

	// Simular criação de usuários
	users := []UserCreatedEvent{
		{
			BaseEvent: events.BaseEvent{Name: "user.created"},
			UserID:    1,
			Username:  "joao_silva",
			Email:     "joao@example.com",
		},
		{
			BaseEvent: events.BaseEvent{Name: "user.created"},
			UserID:    2,
			Username:  "maria_santos",
			Email:     "maria@example.com",
		},
		{
			BaseEvent: events.BaseEvent{Name: "user.created"},
			UserID:    3,
			Username:  "pedro_oliveira",
			Email:     "pedro@example.com",
		},
	}

	appLogger.LogInfo(ctx, "Publishing user creation events", map[string]interface{}{
		"event_type":  "user.created",
		"total_users": len(users),
	})

	// Publicar eventos
	for _, user := range users {
		eventBus.Publish(ctx, &user)
		appLogger.LogInfo(ctx, "Event published", map[string]interface{}{
			"event_type": "user.created",
			"user_id":    user.UserID,
			"username":   user.Username,
		})
		time.Sleep(100 * time.Millisecond)
	}

	appLogger.LogInfo(ctx, "Waiting for event processing", map[string]interface{}{
		"wait_duration": "3s",
	})
	time.Sleep(3 * time.Second)
}

// AsyncPublishingExample demonstrates asynchronous event publishing
func AsyncPublishingExample() {
	ctx := context.Background()
	appLogger := getLogger()

	appLogger.LogInfo(ctx, "Async Publishing Example", map[string]interface{}{
		"example_type": "async_publishing",
		"status":       "starting",
	})

	eventBus := events.NewEventBus()

	// Subscriber para processamento demorado
	slowProcessorChan := make(chan events.Event, 10)
	eventBus.Subscribe("heavy.task", slowProcessorChan)

	go func() {
		for event := range slowProcessorChan {
			appLogger.LogInfo(ctx, "Heavy processing started", map[string]interface{}{
				"service":    "heavy_processor",
				"event_type": event.GetName(),
				"status":     "processing",
			})
			time.Sleep(2 * time.Second)
			appLogger.LogInfo(ctx, "Heavy processing completed", map[string]interface{}{
				"service":    "heavy_processor",
				"event_type": event.GetName(),
				"status":     "completed",
			})
		}
	}()

	event := &events.BaseEvent{Name: "heavy.task"}

	appLogger.LogInfo(ctx, "Publishing heavy task asynchronously", map[string]interface{}{
		"event_type": "heavy.task",
		"mode":       "async",
	})
	eventBus.PublishAsync(ctx, event)

	appLogger.LogInfo(ctx, "Continuing execution without waiting", map[string]interface{}{
		"status": "non_blocking_execution",
	})

	// Simular outras tarefas
	for i := 0; i < 3; i++ {
		appLogger.LogInfo(ctx, "Concurrent task completed", map[string]interface{}{
			"task_number": i + 1,
			"total_tasks": 3,
		})
		time.Sleep(500 * time.Millisecond)
	}

	appLogger.LogInfo(ctx, "Waiting for heavy task completion", map[string]interface{}{
		"wait_duration": "3s",
	})
	time.Sleep(3 * time.Second)
}

// ContextCancellationExample demonstrates context cancellation in EventBus
func ContextCancellationExample() {
	ctx := context.Background()
	appLogger := getLogger()

	appLogger.LogInfo(ctx, "Context Cancellation Example", map[string]interface{}{
		"example_type": "context_cancellation",
		"status":       "starting",
	})

	eventBus := events.NewEventBus()

	subscriber := make(chan events.Event, 10)
	eventBus.Subscribe("cancellable.task", subscriber)

	go func() {
		for event := range subscriber {
			appLogger.LogInfo(ctx, "Processing cancellable task", map[string]interface{}{
				"event_type": event.GetName(),
				"status":     "processing",
			})
		}
	}()

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	event := &events.BaseEvent{Name: "cancellable.task"}

	// Try to publish after context timeout
	time.Sleep(600 * time.Millisecond)

	err := eventBus.Publish(timeoutCtx, event)
	if err != nil {
		appLogger.LogWarn(ctx, "Event publishing cancelled due to context timeout", map[string]interface{}{
			"error":      err.Error(),
			"event_type": "cancellable.task",
		})
	} else {
		appLogger.LogInfo(ctx, "Event published successfully", map[string]interface{}{
			"event_type": "cancellable.task",
		})
	}
}

// BufferOverflowExample demonstrates handling channel buffer overflow
func BufferOverflowExample() {
	ctx := context.Background()
	appLogger := getLogger()

	appLogger.LogInfo(ctx, "Buffer Overflow Example", map[string]interface{}{
		"example_type": "buffer_overflow",
		"status":       "starting",
	})

	eventBus := events.NewEventBus()

	// Small buffer subscriber
	smallBufferChan := make(chan events.Event, 2) // Very small buffer
	eventBus.Subscribe("burst.event", smallBufferChan)

	// Slow consumer
	go func() {
		for event := range smallBufferChan {
			appLogger.LogInfo(ctx, "Slow consumer processing", map[string]interface{}{
				"event_type":      event.GetName(),
				"processing_time": "1s",
			})
			time.Sleep(1 * time.Second) // Slow processing
		}
	}()

	// Fast publisher - will cause buffer overflow
	for i := 0; i < 5; i++ {
		event := &events.BaseEvent{Name: "burst.event"}

		appLogger.LogInfo(ctx, "Publishing burst event", map[string]interface{}{
			"event_number": i + 1,
			"total_events": 5,
		})

		// Use async to avoid blocking
		eventBus.PublishAsync(ctx, event)
		time.Sleep(100 * time.Millisecond)
	}

	appLogger.LogInfo(ctx, "Waiting for burst processing", map[string]interface{}{
		"wait_duration": "6s",
	})
	time.Sleep(6 * time.Second)
}
