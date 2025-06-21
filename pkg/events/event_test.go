package events

import (
	"context"
	"sync"
	"testing"
	"time"
)

type UserCreatedEvent struct {
	*BaseEvent
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

func TestEventBusBasic(t *testing.T) {
	eventBus := NewEventBus(nil)

	var receivedEvent Event
	var mu sync.Mutex

	handler := func(event Event) {
		mu.Lock()
		receivedEvent = event
		mu.Unlock()
	}

	err := eventBus.Subscribe("user.created", handler)
	if err != nil {
		t.Errorf("Error subscribing: %v", err)
	}

	event := &UserCreatedEvent{
		BaseEvent: NewBaseEvent("user.created"),
		UserID:    1,
		Username:  "john_doe",
	}

	err = eventBus.Publish("user.created", event)
	if err != nil {
		t.Errorf("Error publishing event: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	if receivedEvent == nil {
		t.Error("Handler didn't receive event")
	} else {
		userEvent := receivedEvent.(*UserCreatedEvent)
		if userEvent.UserID != 1 || userEvent.Username != "john_doe" {
			t.Errorf("Received event data mismatch")
		}
	}
	mu.Unlock()
}

func TestEventBusSubscribeOnce(t *testing.T) {
	eventBus := NewEventBus(nil)

	callCount := 0
	var mu sync.Mutex

	handler := func(event Event) {
		mu.Lock()
		callCount++
		mu.Unlock()
	}

	err := eventBus.SubscribeOnce("user.created", handler)
	if err != nil {
		t.Errorf("Error subscribing: %v", err)
	}

	event := &UserCreatedEvent{
		BaseEvent: NewBaseEvent("user.created"),
		UserID:    1,
		Username:  "john_doe",
	}

	eventBus.Publish("user.created", event)
	eventBus.Publish("user.created", event)

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	if callCount != 1 {
		t.Errorf("Expected callCount to be 1, got %d", callCount)
	}
	mu.Unlock()
}

func TestEventBusAsync(t *testing.T) {
	eventBus := NewEventBus(nil)

	var wg sync.WaitGroup
	wg.Add(2)

	handler := func(event Event) {
		userEvent := event.(*UserCreatedEvent)
		t.Logf("Async: User created: ID=%d, Username=%s", userEvent.UserID, userEvent.Username)
		wg.Done()
	}

	err := eventBus.SubscribeAsync("user.created", handler, false)
	if err != nil {
		t.Errorf("Error subscribing: %v", err)
	}

	event1 := &UserCreatedEvent{
		BaseEvent: NewBaseEvent("user.created"),
		UserID:    1,
		Username:  "john_doe",
	}

	event2 := &UserCreatedEvent{
		BaseEvent: NewBaseEvent("user.created"),
		UserID:    2,
		Username:  "jane_doe",
	}

	eventBus.PublishAsync("user.created", event1)
	eventBus.PublishAsync("user.created", event2)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Error("Timeout waiting for async events")
	}
}

func TestMultipleSubscribers(t *testing.T) {
	eventBus := NewEventBus(nil)

	received1 := false
	received2 := false
	var mu sync.Mutex

	handler1 := func(event Event) {
		mu.Lock()
		received1 = true
		mu.Unlock()
	}

	handler2 := func(event Event) {
		mu.Lock()
		received2 = true
		mu.Unlock()
	}

	eventBus.Subscribe("test.event", handler1)
	eventBus.Subscribe("test.event", handler2)

	event := NewBaseEvent("test.event")
	err := eventBus.Publish("test.event", event)
	if err != nil {
		t.Errorf("Error publishing event: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	if !received1 || !received2 {
		t.Error("Not all subscribers received the event")
	}
	mu.Unlock()
}

func TestUnsubscribe(t *testing.T) {
	eventBus := NewEventBus(nil)

	callCount := 0
	var mu sync.Mutex

	handler := func(event Event) {
		mu.Lock()
		callCount++
		mu.Unlock()
	}

	eventBus.Subscribe("test.event", handler)

	event := NewBaseEvent("test.event")
	eventBus.Publish("test.event", event)

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	if callCount != 1 {
		t.Errorf("Expected callCount to be 1, got %d", callCount)
	}
	mu.Unlock()

	err := eventBus.Unsubscribe("test.event", handler)
	if err != nil {
		t.Errorf("Error unsubscribing: %v", err)
	}

	eventBus.Publish("test.event", event)
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	if callCount != 1 {
		t.Errorf("Expected callCount to remain 1 after unsubscribe, got %d", callCount)
	}
	mu.Unlock()
}

func TestHasCallback(t *testing.T) {
	eventBus := NewEventBus(nil)

	if eventBus.HasCallback("test.event") {
		t.Error("Should not have callback for non-existent topic")
	}

	handler := func(event Event) {}
	eventBus.Subscribe("test.event", handler)

	if !eventBus.HasCallback("test.event") {
		t.Error("Should have callback for subscribed topic")
	}
}

func TestWaitAsync(t *testing.T) {
	eventBus := NewEventBus(nil)

	var wg sync.WaitGroup
	wg.Add(1)

	handler := func(event Event) {
		time.Sleep(100 * time.Millisecond)
		wg.Done()
	}

	eventBus.SubscribeAsync("test.event", handler, false)

	event := NewBaseEvent("test.event")
	eventBus.PublishAsync("test.event", event)

	start := time.Now()
	eventBus.WaitAsync()
	duration := time.Since(start)

	if duration < 100*time.Millisecond {
		t.Error("WaitAsync should wait for async handlers to complete")
	}
}

func TestChannelEventBus(t *testing.T) {
	channelEventBus := NewChannelEventBus(nil)
	ctx := context.Background()

	subscriber := channelEventBus.SubscribeChannel("user.created", 10)

	event := &UserCreatedEvent{
		BaseEvent: NewBaseEvent("user.created"),
		UserID:    1,
		Username:  "john_doe",
	}

	err := channelEventBus.PublishEvent(ctx, event)
	if err != nil {
		t.Errorf("Error publishing event: %v", err)
	}

	select {
	case receivedEvent := <-subscriber.Channel():
		userEvent := receivedEvent.(*UserCreatedEvent)
		if userEvent.UserID != 1 || userEvent.Username != "john_doe" {
			t.Errorf("Received event data mismatch")
		}
	case <-time.After(time.Second):
		t.Error("Subscriber didn't receive event within timeout")
	}

	subscriber.Close()
}

func TestChannelEventBusAsync(t *testing.T) {
	channelEventBus := NewChannelEventBus(nil)
	ctx := context.Background()

	subscriber := channelEventBus.SubscribeChannel("user.created", 10)

	go func() {
		for event := range subscriber.Channel() {
			userEvent := event.(*UserCreatedEvent)
			t.Logf("Received event: User ID=%d, Username=%s", userEvent.UserID, userEvent.Username)
		}
	}()

	event1 := &UserCreatedEvent{
		BaseEvent: NewBaseEvent("user.created"),
		UserID:    1,
		Username:  "john_doe",
	}

	event2 := &UserCreatedEvent{
		BaseEvent: NewBaseEvent("user.created"),
		UserID:    2,
		Username:  "jane_doe",
	}

	channelEventBus.PublishEventAsync(ctx, event1)
	channelEventBus.PublishEventAsync(ctx, event2)

	time.Sleep(200 * time.Millisecond)
	subscriber.Close()
}
