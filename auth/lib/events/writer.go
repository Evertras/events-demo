package events

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"

	"github.com/Evertras/events-demo/auth/lib/events/authevents"
	"github.com/Evertras/events-demo/shared/stream"
)

// Writer can write to an event stream with predefined event types
type Writer interface {
	// PostRegisteredEvent posts a UserRegistered event to the stream
	PostRegisteredEvent(ctx context.Context, ev *authevents.UserRegistered) error

	// Close gives underlying resources a chance to flush and close gracefully
	Close() error
}

type writer struct {
	streamWriter stream.Writer
}

func NewWriter(streamWriter stream.Writer) Writer {
	return &writer{
		streamWriter: streamWriter,
	}
}

func (w *writer) PostRegisteredEvent(ctx context.Context, ev *authevents.UserRegistered) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostRegisteredEvent")
	defer span.Finish()

	return w.streamWriter.Write(ctx, []byte(ev.ID), EventIDUserRegistered, ev)
}

func (w *writer) Close() error {
	return w.streamWriter.Close()
}
