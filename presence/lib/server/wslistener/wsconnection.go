package wslistener

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"

	"github.com/Evertras/events-demo/presence/lib/friendlist"
)

type wsConnection struct {
	conn *websocket.Conn
	done chan interface{}
	id   string
}

func newConnection(conn *websocket.Conn, id string) *wsConnection {
	return &wsConnection{
		conn: conn,
		done: make(chan interface{}),
		id:   id,
	}
}

func (c *wsConnection) Close() {
	close(c.done)
}

func (c *wsConnection) Done() chan interface{} {
	return c.done
}

func (c *wsConnection) GetID() string {
	return c.id
}

func onlineText(isOnline bool) string {
	if isOnline {
		return "online"
	}

	return "offline"
}

func (c *wsConnection) SendFriendStatusNotification(notification friendlist.FriendStatus) error {
	msg := fmt.Sprintf("%s is %s", notification.ID, onlineText(notification.Online))

	log.Println("Notifying", c.id, "that", msg)

	c.conn.WriteMessage(websocket.TextMessage, []byte(msg))

	return nil
}
