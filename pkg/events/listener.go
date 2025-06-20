package events

import "context"

type Listener interface {
	Handle(ctx context.Context, event Event) error
}

type ListenerFunc func(ctx context.Context, event Event) error

func (f ListenerFunc) Handle(ctx context.Context, event Event) error {
	return f(ctx, event)
}
