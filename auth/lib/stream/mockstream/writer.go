package mockstream

import (
	"context"

	"github.com/Evertras/events-demo/auth/lib/stream"
	"github.com/Evertras/events-demo/auth/lib/stream/authevents"
)

type MockStreamWriter struct {
	Sent []interface{}
}

var _ stream.Writer = &MockStreamWriter{}

func NewWriter() *MockStreamWriter {
	return &MockStreamWriter{
		Sent: make([]interface{}, 0),
	}
}

func (s *MockStreamWriter) PostRegisteredEvent(ctx context.Context, ev *authevents.UserRegistered) error {
	s.Sent = append(s.Sent, ev)

	return nil
}

func (s *MockStreamWriter) GetNextRegisteredEventFor(ctx context.Context, id string) (*authevents.UserRegistered, error) {
	return nil, nil
}

func (s *MockStreamWriter) Close() error {
	return nil
}
