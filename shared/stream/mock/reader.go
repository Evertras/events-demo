package mock

import (
	"context"

	"github.com/Evertras/events-demo/shared/stream"
)

type MockReceivedEvent struct {
	ID   stream.EventID
	Data []byte
}

type MockStreamReader struct {
	Receive  chan MockReceivedEvent
	Handlers map[stream.EventID]stream.StreamEventHandler
}

var _ stream.Reader = &MockStreamReader{}

func NewReader() *MockStreamReader {
	return &MockStreamReader{
		Receive:  make(chan MockReceivedEvent, 100),
		Handlers: make(map[stream.EventID]stream.StreamEventHandler),
	}
}

func (s *MockStreamReader) RegisterHandler(event stream.EventID, handler stream.StreamEventHandler) error {
	s.Handlers[event] = handler
	return nil
}

func (s *MockStreamReader) Listen(ctx context.Context) error {
	for {
		select {
		case ev := <-s.Receive:
			if handler, ok := s.Handlers[ev.ID]; ok {
				handler(ctx, ev.Data)
			}

		case <-ctx.Done():
			return nil
		}
	}
}
