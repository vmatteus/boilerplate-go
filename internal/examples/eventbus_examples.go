package examples

import (
	"context"
	"time"

	"github.com/your-org/boilerplate-go/pkg/events"
)

type UserCreatedEvent struct {
	*events.BaseEvent
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// NewEventBusExample demonstrates the EventBus with multiple subscribers
func NewEventBusExample() {
	ctx := context.Background()
	appLogger := getLogger()

	appLogger.LogInfo(ctx, "EventBus Handler-Based Example", map[string]interface{}{
		"example_type": "eventbus_handlers",
		"status":       "starting",
	})

	eventBus := events.NewEventBus(nil)

	emailHandler := func(event events.Event) {
		user := event.(*UserCreatedEvent)
		appLogger.LogInfo(ctx, "Email service processing", map[string]interface{}{
			"service":  "email",
			"action":   "send_welcome_email",
			"username": user.Username,
			"email":    user.Email,
		})
		time.Sleep(500 * time.Millisecond)
	}

	analyticsHandler := func(event events.Event) {
		user := event.(*UserCreatedEvent)
		appLogger.LogInfo(ctx, "Analytics service processing", map[string]interface{}{
			"service": "analytics",
			"action":  "record_user_metrics",
			"user_id": user.UserID,
		})
		time.Sleep(200 * time.Millisecond)
	}

	notificationHandler := func(event events.Event) {
		user := event.(*UserCreatedEvent)
		appLogger.LogInfo(ctx, "Notification service processing", map[string]interface{}{
			"service":  "notifications",
			"action":   "create_welcome_notification",
			"username": user.Username,
		})
		time.Sleep(300 * time.Millisecond)
	}

	eventBus.Subscribe("user.created", emailHandler)
	eventBus.Subscribe("user.created", analyticsHandler)
	eventBus.Subscribe("user.created", notificationHandler)

	users := []UserCreatedEvent{
		{
			BaseEvent: events.NewBaseEvent("user.created"),
			UserID:    1,
			Username:  "joao_silva",
			Email:     "joao@example.com",
		},
		{
			BaseEvent: events.NewBaseEvent("user.created"),
			UserID:    2,
			Username:  "maria_santos",
			Email:     "maria@example.com",
		},
		{
			BaseEvent: events.NewBaseEvent("user.created"),
			UserID:    3,
			Username:  "pedro_oliveira",
			Email:     "pedro@example.com",
		},
	}

	appLogger.LogInfo(ctx, "Publishing user creation events", map[string]interface{}{
		"event_type":  "user.created",
		"total_users": len(users),
	})

	for _, user := range users {
		eventBus.Publish("user.created", &user)
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

	eventBus := events.NewEventBus(nil)

	slowProcessorHandler := func(event events.Event) {
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

	eventBus.SubscribeAsync("heavy.task", slowProcessorHandler, false)

	event := events.NewBaseEvent("heavy.task")

	appLogger.LogInfo(ctx, "Publishing heavy task asynchronously", map[string]interface{}{
		"event_type": "heavy.task",
		"mode":       "async",
	})
	eventBus.PublishAsync("heavy.task", event)

	appLogger.LogInfo(ctx, "Continuing execution without waiting", map[string]interface{}{
		"status": "non_blocking_execution",
	})

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
	eventBus.WaitAsync()
}

// ContextCancellationExample demonstrates context cancellation in EventBus
func ContextCancellationExample() {
	ctx := context.Background()
	appLogger := getLogger()

	appLogger.LogInfo(ctx, "Context Cancellation Example", map[string]interface{}{
		"example_type": "context_cancellation",
		"status":       "starting",
	})

	eventBus := events.NewEventBus(nil)

	handler := func(event events.Event) {
		appLogger.LogInfo(ctx, "Processing cancellable task", map[string]interface{}{
			"event_type": event.GetName(),
			"status":     "processing",
		})
	}

	eventBus.Subscribe("cancellable.task", handler)

	event := events.NewBaseEvent("cancellable.task")

	appLogger.LogInfo(ctx, "Publishing event", map[string]interface{}{
		"event_type": "cancellable.task",
	})

	err := eventBus.Publish("cancellable.task", event)
	if err != nil {
		appLogger.LogWarn(ctx, "Event publishing failed", map[string]interface{}{
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

	eventBus := events.NewEventBus(nil)

	slowHandler := func(event events.Event) {
		appLogger.LogInfo(ctx, "Slow consumer processing", map[string]interface{}{
			"event_type":      event.GetName(),
			"processing_time": "1s",
		})
		time.Sleep(1 * time.Second)
	}

	eventBus.SubscribeAsync("burst.event", slowHandler, false)

	for i := 0; i < 5; i++ {
		event := events.NewBaseEvent("burst.event")

		appLogger.LogInfo(ctx, "Publishing burst event", map[string]interface{}{
			"event_number": i + 1,
			"total_events": 5,
		})

		eventBus.PublishAsync("burst.event", event)
		time.Sleep(100 * time.Millisecond)
	}

	appLogger.LogInfo(ctx, "Waiting for burst processing", map[string]interface{}{
		"wait_duration": "6s",
	})
	eventBus.WaitAsync()
	time.Sleep(1 * time.Second)
}

func SubscribeOnceExample() {
	ctx := context.Background()
	appLogger := getLogger()

	appLogger.LogInfo(ctx, "Subscribe Once Example", map[string]interface{}{
		"example_type": "subscribe_once",
		"status":       "starting",
	})

	eventBus := events.NewEventBus(nil)

	oneTimeHandler := func(event events.Event) {
		appLogger.LogInfo(ctx, "One-time handler executed", map[string]interface{}{
			"event_type": event.GetName(),
			"event_id":   event.GetID(),
		})
	}

	eventBus.SubscribeOnce("one.time.event", oneTimeHandler)

	event := events.NewBaseEvent("one.time.event")

	appLogger.LogInfo(ctx, "Publishing event first time", map[string]interface{}{
		"event_type": "one.time.event",
	})
	eventBus.Publish("one.time.event", event)

	appLogger.LogInfo(ctx, "Publishing event second time (should not trigger handler)", map[string]interface{}{
		"event_type": "one.time.event",
	})
	eventBus.Publish("one.time.event", event)

	time.Sleep(500 * time.Millisecond)
}

func ChannelEventBusExample() {
	ctx := context.Background()
	appLogger := getLogger()

	appLogger.LogInfo(ctx, "Channel EventBus Example", map[string]interface{}{
		"example_type": "channel_eventbus",
		"status":       "starting",
	})

	channelEventBus := events.NewChannelEventBus(nil)

	emailSubscriber := channelEventBus.SubscribeChannel("user.created", 10)
	analyticsSubscriber := channelEventBus.SubscribeChannel("user.created", 10)

	go func() {
		for event := range emailSubscriber.Channel() {
			user := event.(*UserCreatedEvent)
			appLogger.LogInfo(ctx, "Email service processing via channel", map[string]interface{}{
				"service":  "email",
				"action":   "send_welcome_email",
				"username": user.Username,
				"email":    user.Email,
			})
		}
	}()

	go func() {
		for event := range analyticsSubscriber.Channel() {
			user := event.(*UserCreatedEvent)
			appLogger.LogInfo(ctx, "Analytics service processing via channel", map[string]interface{}{
				"service": "analytics",
				"action":  "record_user_metrics",
				"user_id": user.UserID,
			})
		}
	}()

	users := []UserCreatedEvent{
		{
			BaseEvent: events.NewBaseEvent("user.created"),
			UserID:    1,
			Username:  "joao_silva",
			Email:     "joao@example.com",
		},
		{
			BaseEvent: events.NewBaseEvent("user.created"),
			UserID:    2,
			Username:  "maria_santos",
			Email:     "maria@example.com",
		},
	}

	for _, user := range users {
		channelEventBus.PublishEvent(ctx, &user)
		appLogger.LogInfo(ctx, "Event published via channel", map[string]interface{}{
			"event_type": "user.created",
			"user_id":    user.UserID,
			"username":   user.Username,
		})
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(2 * time.Second)

	emailSubscriber.Close()
	analyticsSubscriber.Close()
	channelEventBus.Close()
}
