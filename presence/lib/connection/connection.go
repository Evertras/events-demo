package connection

import (
	"github.com/Evertras/events-demo/presence/lib/friendlist"
)

// Connection is a player connection
type Connection interface {
	// Listen blocks until the connection disconnects
	Listen() error

	SendFriendStatusNotification(notification friendlist.FriendStatus) error
}
