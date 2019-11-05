package connection

import (
	"github.com/Evertras/events-demo/presence/lib/friendlist"
)

// Connection is a verified player connection that has authenticated itself
// somehow and is identified with a player ID
type Connection interface {
	// Close closes the connection
	Close()

	// Done returns a channel that closes when the connection disconnects
	Done() chan interface{}

	// GetID returns the player ID of the connection
	GetID() string

	SendFriendStatusNotification(notification friendlist.FriendStatus) error
}
