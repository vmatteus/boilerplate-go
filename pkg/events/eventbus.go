package events

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

type EventBus interface {
	Subscribe(topic string, handler interface{}) error
	SubscribeOnce(topic string, handler interface{}) error
	SubscribeAsync(topic string, handler interface{}, transactional bool) error
	SubscribeOnceAsync(topic string, handler interface{}) error
	Unsubscribe(topic string, handler interface{}) error
	Publish(topic string, args ...interface{}) error
	PublishAsync(topic string, args ...interface{})
	HasCallback(topic string) bool
	WaitAsync()
	Close() error
}

type eventBus struct {
	handlers map[string][]*eventHandler
	mu       sync.RWMutex
	wg       sync.WaitGroup
	closed   bool
	closeMu  sync.RWMutex
}

type eventHandler struct {
	callBack      reflect.Value
	once          bool
	async         bool
	transactional bool
}

type EventBusConfig struct {
	DefaultBufferSize int
	DefaultTimeout    time.Duration
}

func DefaultConfig() *EventBusConfig {
	return &EventBusConfig{
		DefaultBufferSize: 100,
		DefaultTimeout:    30 * time.Second,
	}
}

func NewEventBus(config *EventBusConfig) EventBus {
	if config == nil {
		config = DefaultConfig()
	}

	return &eventBus{
		handlers: make(map[string][]*eventHandler),
	}
}

func (bus *eventBus) Subscribe(topic string, fn interface{}) error {
	return bus.subscribe(topic, fn, false, false, false)
}

func (bus *eventBus) SubscribeOnce(topic string, fn interface{}) error {
	return bus.subscribe(topic, fn, true, false, false)
}

func (bus *eventBus) SubscribeAsync(topic string, fn interface{}, transactional bool) error {
	return bus.subscribe(topic, fn, false, true, transactional)
}

func (bus *eventBus) SubscribeOnceAsync(topic string, fn interface{}) error {
	return bus.subscribe(topic, fn, true, true, false)
}

func (bus *eventBus) subscribe(topic string, fn interface{}, once, async, transactional bool) error {
	bus.closeMu.RLock()
	if bus.closed {
		bus.closeMu.RUnlock()
		return fmt.Errorf("eventbus is closed")
	}
	bus.closeMu.RUnlock()

	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("%s is not of type reflect.Func", reflect.TypeOf(fn))
	}

	bus.mu.Lock()
	defer bus.mu.Unlock()

	handler := &eventHandler{
		callBack:      reflect.ValueOf(fn),
		once:          once,
		async:         async,
		transactional: transactional,
	}

	bus.handlers[topic] = append(bus.handlers[topic], handler)
	return nil
}

func (bus *eventBus) Unsubscribe(topic string, fn interface{}) error {
	bus.closeMu.RLock()
	if bus.closed {
		bus.closeMu.RUnlock()
		return fmt.Errorf("eventbus is closed")
	}
	bus.closeMu.RUnlock()

	bus.mu.Lock()
	defer bus.mu.Unlock()

	if _, ok := bus.handlers[topic]; !ok {
		return fmt.Errorf("topic %s doesn't exist", topic)
	}

	rv := reflect.ValueOf(fn)
	for i, handler := range bus.handlers[topic] {
		if handler.callBack == rv {
			bus.handlers[topic] = append(bus.handlers[topic][:i], bus.handlers[topic][i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("handler not found for topic %s", topic)
}

func (bus *eventBus) Publish(topic string, args ...interface{}) error {
	bus.closeMu.RLock()
	if bus.closed {
		bus.closeMu.RUnlock()
		return fmt.Errorf("eventbus is closed")
	}
	bus.closeMu.RUnlock()

	bus.mu.RLock()
	handlers := make([]*eventHandler, len(bus.handlers[topic]))
	copy(handlers, bus.handlers[topic])
	bus.mu.RUnlock()

	if len(handlers) > 0 {
		for _, handler := range handlers {
			if handler.async {
				bus.wg.Add(1)
				go bus.executeHandler(handler, args...)
			} else {
				bus.executeHandler(handler, args...)
			}
		}
	}
	return nil
}

func (bus *eventBus) PublishAsync(topic string, args ...interface{}) {
	bus.closeMu.RLock()
	if bus.closed {
		bus.closeMu.RUnlock()
		return
	}
	bus.closeMu.RUnlock()

	bus.mu.RLock()
	handlers := make([]*eventHandler, len(bus.handlers[topic]))
	copy(handlers, bus.handlers[topic])
	bus.mu.RUnlock()

	if len(handlers) > 0 {
		for _, handler := range handlers {
			bus.wg.Add(1)
			go bus.executeHandler(handler, args...)
		}
	}
}

func (bus *eventBus) executeHandler(handler *eventHandler, args ...interface{}) {
	if handler.async {
		defer bus.wg.Done()
	}

	passedArguments := bus.setUpPublish(handler.callBack, args...)
	handler.callBack.Call(passedArguments)

	if handler.once {
		if handler.async {
			go bus.removeHandler(handler)
		} else {
			bus.removeHandler(handler)
		}
	}
}

func (bus *eventBus) removeHandler(handler *eventHandler) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	for topic, handlers := range bus.handlers {
		for i, h := range handlers {
			if h == handler {
				bus.handlers[topic] = append(handlers[:i], handlers[i+1:]...)
				return
			}
		}
	}
}

func (bus *eventBus) setUpPublish(function reflect.Value, args ...interface{}) []reflect.Value {
	funcType := function.Type()
	passedArguments := make([]reflect.Value, len(args))

	for i, v := range args {
		if v == nil {
			passedArguments[i] = reflect.New(funcType.In(i)).Elem()
		} else {
			passedArguments[i] = reflect.ValueOf(v)
		}
	}

	return passedArguments
}

func (bus *eventBus) HasCallback(topic string) bool {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	_, ok := bus.handlers[topic]
	return ok && len(bus.handlers[topic]) > 0
}

func (bus *eventBus) WaitAsync() {
	bus.wg.Wait()
}

func (bus *eventBus) Close() error {
	bus.closeMu.Lock()
	defer bus.closeMu.Unlock()

	if bus.closed {
		return nil
	}

	bus.closed = true
	bus.WaitAsync()

	bus.mu.Lock()
	bus.handlers = make(map[string][]*eventHandler)
	bus.mu.Unlock()

	return nil
}
