package mockeventstream

import (
	"github.com/Evertras/events-demo/auth/lib/eventstream"
	"github.com/Evertras/events-demo/auth/lib/eventstream/events"
)

type MockEventStream struct {
	Sent []interface{}
}

var _ eventstream.EventStream = &MockEventStream{}

func New() *MockEventStream {
	return &MockEventStream{
		Sent: make([]interface{}, 0),
	}
}

func (s *MockEventStream) PostRegisteredEvent(ev *events.UserRegistered) error {
	s.Sent = append(s.Sent, ev)

	return nil
}
