package wslistener

import (
	"context"
	"log"
	"net/http"

	"github.com/Evertras/events-demo/presence/lib/connection"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type WsListener struct {
	connections chan connection.Connection
	server      http.Server
}

func New(addr string, route string) *WsListener {
	mux := http.NewServeMux()
	connections := make(chan connection.Connection, 100)

	mux.HandleFunc(route, upgradeHandler(connections))

	return &WsListener{
		server: http.Server{
			Addr:    addr,
			Handler: mux,
		},
		connections: connections,
	}
}

func upgradeHandler(conns chan connection.Connection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println("Unexpected error when trying to upgrade:", err)
			return
		}
		defer c.Close()

		_, raw, err := c.ReadMessage()

		id := string(raw)

		if err != nil {
			log.Println("Failed to read:", err)
			return
		}

		log.Println("Connecting as", id)

		conn := newConnection(c, id)

		defer conn.Close()

		conns <- conn

		for {
			mt, msg, err := c.ReadMessage()

			if err != nil {
				log.Println("Failed to read:", err)
				break
			}

			log.Println("Received:", mt, string(msg))
		}
	}
}

func (l *WsListener) Listen(ctx context.Context) error {
	log.Println("Listening...")
	return l.server.ListenAndServe()
}

func (l *WsListener) IdentifiedConnections() <-chan connection.Connection {
	return l.connections
}
