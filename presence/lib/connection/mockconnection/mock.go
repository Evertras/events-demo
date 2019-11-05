package mockconnection

import (
	"sync"

	"github.com/Evertras/events-demo/presence/lib/friendlist"
)

type MockConnection struct {
	m sync.Mutex

	id                    string
	closed                chan interface{}
	receivedNotifications []friendlist.FriendStatus
}

func New(id string) *MockConnection {
	return &MockConnection{
		id:                    id,
		closed:                make(chan interface{}),
		receivedNotifications: make([]friendlist.FriendStatus, 0),
	}
}

////////////////////////////////////////////////////////////////////////////////
// Interface stuff
func (c *MockConnection) Done() chan interface{} {
	return c.closed
}

func (c *MockConnection) GetID() string {
	return c.id
}

func (c *MockConnection) SendFriendStatusNotification(n friendlist.FriendStatus) error {
	c.m.Lock()
	c.receivedNotifications = append(c.receivedNotifications, n)
	c.m.Unlock()

	return nil
}

func (c *MockConnection) Close() {
	close(c.closed)
}

////////////////////////////////////////////////////////////////////////////////
// Mock stuff

func (c *MockConnection) GetReceivedFriendStatusNotifications() []friendlist.FriendStatus {
	c.m.Lock()
	defer c.m.Unlock()
	return c.receivedNotifications
}
