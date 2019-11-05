package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8111", "HTTP server address")
var playerID = flag.String("id", "A", "Player ID")

func main() {
	flag.Parse()
	done := make(chan struct{})
	interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:8111",
		Path:   "/",
	}

	log.Printf("Connecting to %s as %q", u.String(), *playerID)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatal("Dial:", err)
	}
	defer c.Close()

	go func() {
		defer close(done)

		for {
			_, message, err := c.ReadMessage()

			if err != nil {
				log.Println("Read error:", err)
				return
			}

			log.Println("Received:", string(message))
		}
	}()

	err = c.WriteMessage(websocket.TextMessage, []byte(*playerID))

	if err != nil {
		log.Println("Write error:", err)
		return
	}

	for {
		select {
		case <-done:
			return

		case <-interrupt:
			log.Println("Interrupted")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

			if err != nil {
				log.Println("Close error:", err)
				return
			}

			select {
			case <-done:
			case <-time.After(time.Millisecond * 200):
			}

			return
		}
	}
}
