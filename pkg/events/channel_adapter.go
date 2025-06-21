package events

import (
	"context"
	"sync"
)

type ChannelSubscriber struct {
	channel chan Event
	ctx     context.Context
	cancel  context.CancelFunc
	mu      sync.RWMutex
	closed  bool
}

type ChannelEventBus struct {
	EventBus
	subscribers map[string][]*ChannelSubscriber
	mu          sync.RWMutex
}

func NewChannelEventBus(config *EventBusConfig) *ChannelEventBus {
	return &ChannelEventBus{
		EventBus:    NewEventBus(config),
		subscribers: make(map[string][]*ChannelSubscriber),
	}
}

func (ceb *ChannelEventBus) SubscribeChannel(topic string, bufferSize int) *ChannelSubscriber {
	if bufferSize <= 0 {
		bufferSize = 100
	}

	subscriber := &ChannelSubscriber{
		channel: make(chan Event, bufferSize),
	}
	subscriber.ctx, subscriber.cancel = context.WithCancel(context.Background())

	ceb.mu.Lock()
	ceb.subscribers[topic] = append(ceb.subscribers[topic], subscriber)
	ceb.mu.Unlock()

	return subscriber
}

func (ceb *ChannelEventBus) PublishEvent(ctx context.Context, event Event) error {
	ceb.mu.RLock()
	subscribers := ceb.subscribers[event.GetName()]
	ceb.mu.RUnlock()

	for _, subscriber := range subscribers {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-subscriber.ctx.Done():
			continue
		case subscriber.channel <- event:
		default:
		}
	}

	return nil
}

func (ceb *ChannelEventBus) PublishEventAsync(ctx context.Context, event Event) {
	ceb.mu.RLock()
	subscribers := ceb.subscribers[event.GetName()]
	ceb.mu.RUnlock()

	for _, subscriber := range subscribers {
		go func(sub *ChannelSubscriber) {
			select {
			case <-ctx.Done():
				return
			case <-sub.ctx.Done():
				return
			case sub.channel <- event:
			default:
			}
		}(subscriber)
	}
}

func (ceb *ChannelEventBus) UnsubscribeChannel(topic string, subscriber *ChannelSubscriber) {
	ceb.mu.Lock()
	defer ceb.mu.Unlock()

	if subscribers, exists := ceb.subscribers[topic]; exists {
		for i, sub := range subscribers {
			if sub == subscriber {
				ceb.subscribers[topic] = append(subscribers[:i], subscribers[i+1:]...)
				subscriber.Close()
				break
			}
		}
	}
}

func (ceb *ChannelEventBus) Close() error {
	ceb.mu.Lock()
	defer ceb.mu.Unlock()

	for _, subscribers := range ceb.subscribers {
		for _, subscriber := range subscribers {
			subscriber.Close()
		}
	}

	ceb.subscribers = make(map[string][]*ChannelSubscriber)
	return ceb.EventBus.Close()
}

func (cs *ChannelSubscriber) Channel() <-chan Event {
	return cs.channel
}

func (cs *ChannelSubscriber) Close() {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if !cs.closed {
		cs.closed = true
		cs.cancel()
		close(cs.channel)
	}
}

func (cs *ChannelSubscriber) IsClosed() bool {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.closed
}

func (cs *ChannelSubscriber) Context() context.Context {
	return cs.ctx
}
