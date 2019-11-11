package mock

import (
	"context"
	"log"

	"github.com/Evertras/events-demo/shared/stream"
)

type MockReceivedEvent struct {
	ID   stream.EventID
	Data []byte
}

type MockStreamReader struct {
	Receive  chan MockReceivedEvent
	Handlers map[stream.EventID]stream.StreamEventHandler
	Logger   *log.Logger
}

var _ stream.Reader = &MockStreamReader{}

type MockStreamReaderOpts struct {
	Logger *log.Logger
}

func NewReader(opts MockStreamReaderOpts) *MockStreamReader {
	return &MockStreamReader{
		Receive:  make(chan MockReceivedEvent, 100),
		Handlers: make(map[stream.EventID]stream.StreamEventHandler),
		Logger:   opts.Logger,
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
				err := handler(ctx, ev.Data)

				if err != nil && s.Logger != nil {
					s.Logger.Println(err)
				}
			}

		case <-ctx.Done():
			return nil
		}
	}
}
