package events

import (
	"context"
	"sync"
)

type EventBus struct {
	subscribers map[string][]chan<- Event
	mu          sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan<- Event),
	}
}

func (eb *EventBus) Subscribe(eventType string, subscriber chan<- Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.subscribers[eventType] = append(eb.subscribers[eventType], subscriber)
}

func (eb *EventBus) Publish(ctx context.Context, event Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	eb.mu.RLock()
	subscribers := eb.subscribers[event.GetName()]
	eb.mu.RUnlock()

	for _, subscriber := range subscribers {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case subscriber <- event:
		}
	}

	return nil
}

func (eb *EventBus) PublishAsync(ctx context.Context, event Event) {
	eb.mu.RLock()
	subscribers := eb.subscribers[event.GetName()]
	eb.mu.RUnlock()

	for _, subscriber := range subscribers {
		go func(sub chan<- Event) {
			select {
			case <-ctx.Done():
				return
			case sub <- event:
			default:
			}
		}(subscriber)
	}
}
