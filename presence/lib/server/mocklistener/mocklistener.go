package mocklistener

import (
	"context"

	"github.com/Evertras/events-demo/presence/lib/connection"
)

type MockListener struct {
	connections chan connection.Connection
}

func New() *MockListener {
	return &MockListener{
		connections: make(chan connection.Connection, 10),
	}
}

///////////////////////////////////////////////////////////////////////////////
// Interface stuff

func (l *MockListener) Listen(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (l *MockListener) IdentifiedConnections() <-chan connection.Connection {
	return l.connections
}

///////////////////////////////////////////////////////////////////////////////
// Mock stuff

func (l *MockListener) AddConnection(c connection.Connection) {
	l.connections <- c
}
