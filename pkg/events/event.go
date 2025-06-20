package events

type Event interface {
	GetName() string
}

type BaseEvent struct {
	Name string
}

func (e *BaseEvent) GetName() string {
	return e.Name
}
