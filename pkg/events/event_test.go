package events

import (
	"context"
	"sync"
	"testing"
	"time"
)

type UserCreatedEvent struct {
	BaseEvent
	UserID   int
	Username string
}

func TestEventBusBasic(t *testing.T) {
	eventBus := NewEventBus()
	ctx := context.Background()

	subscriber := make(chan Event, 10)
	eventBus.Subscribe("user.created", subscriber)

	event := &UserCreatedEvent{
		BaseEvent: BaseEvent{Name: "user.created"},
		UserID:    1,
		Username:  "john_doe",
	}

	err := eventBus.Publish(ctx, event)
	if err != nil {
		t.Errorf("Error publishing event: %v", err)
	}

	select {
	case receivedEvent := <-subscriber:
		userEvent := receivedEvent.(*UserCreatedEvent)
		if userEvent.UserID != 1 || userEvent.Username != "john_doe" {
			t.Errorf("Received event data mismatch")
		}
	case <-time.After(time.Second):
		t.Error("Subscriber didn't receive event within timeout")
	}
}

func TestEventBusAsync(t *testing.T) {
	eventBus := NewEventBus()
	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(2)

	subscriber := make(chan Event, 10)
	eventBus.Subscribe("user.created", subscriber)

	go func() {
		for event := range subscriber {
			userEvent := event.(*UserCreatedEvent)
			t.Logf("Async: User created: ID=%d, Username=%s", userEvent.UserID, userEvent.Username)
			wg.Done()
		}
	}()

	event1 := &UserCreatedEvent{
		BaseEvent: BaseEvent{Name: "user.created"},
		UserID:    1,
		Username:  "john_doe",
	}

	event2 := &UserCreatedEvent{
		BaseEvent: BaseEvent{Name: "user.created"},
		UserID:    2,
		Username:  "jane_doe",
	}

	eventBus.PublishAsync(ctx, event1)
	eventBus.PublishAsync(ctx, event2)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(3 * time.Second):
		t.Error("Timeout waiting for async events")
	}
}

func TestMultipleSubscribers(t *testing.T) {
	eventBus := NewEventBus()
	ctx := context.Background()

	subscriber1 := make(chan Event, 10)
	subscriber2 := make(chan Event, 10)

	eventBus.Subscribe("test.event", subscriber1)
	eventBus.Subscribe("test.event", subscriber2)

	event := &BaseEvent{Name: "test.event"}
	err := eventBus.Publish(ctx, event)
	if err != nil {
		t.Errorf("Error publishing event: %v", err)
	}

	select {
	case receivedEvent := <-subscriber1:
		if receivedEvent.GetName() != "test.event" {
			t.Errorf("Expected event name 'test.event', got '%s'", receivedEvent.GetName())
		}
	case <-time.After(time.Second):
		t.Error("Subscriber1 didn't receive event within timeout")
	}

	select {
	case receivedEvent := <-subscriber2:
		if receivedEvent.GetName() != "test.event" {
			t.Errorf("Expected event name 'test.event', got '%s'", receivedEvent.GetName())
		}
	case <-time.After(time.Second):
		t.Error("Subscriber2 didn't receive event within timeout")
	}
}

func TestContextCancellation(t *testing.T) {
	eventBus := NewEventBus()
	ctx, cancel := context.WithCancel(context.Background())

	subscriber := make(chan Event, 10)
	eventBus.Subscribe("test.event", subscriber)

	cancel()

	event := &BaseEvent{Name: "test.event"}
	err := eventBus.Publish(ctx, event)

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got: %v", err)
	}
}

func TestAsyncPublishing(t *testing.T) {
	eventBus := NewEventBus()
	ctx := context.Background()

	subscriber := make(chan Event, 10)
	eventBus.Subscribe("test.event", subscriber)

	event := &BaseEvent{Name: "test.event"}
	eventBus.PublishAsync(ctx, event)

	select {
	case receivedEvent := <-subscriber:
		if receivedEvent.GetName() != "test.event" {
			t.Errorf("Expected event name 'test.event', got '%s'", receivedEvent.GetName())
		}
	case <-time.After(time.Second):
		t.Error("Subscriber didn't receive event within timeout")
	}
}
