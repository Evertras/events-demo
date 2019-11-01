package mockconnection

import (
	"github.com/Evertras/events-demo/presence/lib/connection"
	"github.com/Evertras/events-demo/presence/lib/friendlist"
)

type mockConnection struct {
	closed chan bool
}

func NewMockConnection() connection.Connection {
	return &mockConnection{
		closed: make(chan bool, 1),
	}
}

func (c *mockConnection) Listen() error {
	<-c.closed
	return nil
}

func (c *mockConnection) SendFriendStatusNotification(friendlist.FriendStatus) error {
	return nil
}

func (c *mockConnection) Close() {
	c.closed <- true
}
