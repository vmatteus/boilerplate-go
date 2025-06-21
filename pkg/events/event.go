package events

import (
	"time"
)

type Event interface {
	GetName() string
	GetTimestamp() time.Time
	GetID() string
}

type BaseEvent struct {
	Name      string    `json:"name"`
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

func NewBaseEvent(name string) *BaseEvent {
	return &BaseEvent{
		Name:      name,
		ID:        generateEventID(),
		Timestamp: time.Now(),
	}
}

func (e *BaseEvent) GetName() string {
	return e.Name
}

func (e *BaseEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

func (e *BaseEvent) GetID() string {
	return e.ID
}

func generateEventID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
