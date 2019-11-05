package server

import (
	"context"
	"log"
	"sync"

	"github.com/Evertras/events-demo/presence/lib/connection"
	"github.com/Evertras/events-demo/presence/lib/db"
	"github.com/Evertras/events-demo/presence/lib/friendlist"
)

type Listener interface {
	Listen(ctx context.Context) error

	IdentifiedConnections() <-chan connection.Connection
}

type Server interface {
	// Run blocks until an error occurs or the context closes
	Run(ctx context.Context) error
}

type server struct {
	m sync.Mutex

	connections map[string]connection.Connection
	db          db.Db
	fl          friendlist.FriendList
	listener    Listener
}

func New(listener Listener, d db.Db, fl friendlist.FriendList) Server {
	return &server{
		connections: make(map[string]connection.Connection),
		db:          d,
		fl:          fl,
		listener:    listener,
	}
}

func (s *server) Run(ctx context.Context) error {
	identifiedConnections := s.listener.IdentifiedConnections()

	go s.listener.Listen(ctx)

	for {
		select {
		case c, open := <-identifiedConnections:
			if !open {
				return nil
			}

			if err := s.addConnection(ctx, c); err != nil {
				return err
			}

		case <-ctx.Done():
			return nil
		}
	}
}

func (s *server) addConnection(ctx context.Context, c connection.Connection) error {
	id := c.GetID()

	friends, err := s.db.GetFriendList(ctx, id)

	if err != nil {
		return err
	}

	err = s.db.SendNotification(ctx, db.PresenceChangedEvent{
		PlayerID:  id,
		NotifyIDs: friends,
		Online:    true,
	})

	if err != nil {
		return err
	}

	s.m.Lock()
	if existing, exists := s.connections[id]; exists {
		existing.Close()
	}
	s.connections[id] = c
	s.m.Unlock()

	// Feed friend list notifications to the connection
	go func() {
		notifications, err := s.fl.Subscribe(ctx, id)

		if err != nil {
			// TODO: not this
			log.Println(err)
			return
		}

		for {
			select {
			case n, open := <-notifications:
				if !open {
					return
				}

				c.SendFriendStatusNotification(n)

			case <-ctx.Done():
				return
			}
		}
	}()

	// Clean up on close
	go func() {
		<-c.Done()

		s.fl.Unsubscribe(ctx, id)

		s.m.Lock()
		delete(s.connections, id)
		s.m.Unlock()

		// friends list maybe changed since the start
		friends, err := s.db.GetFriendList(ctx, id)

		if err != nil {
			log.Println("Error getting friends on close cleanup:", err)
		}

		err = s.db.SendNotification(ctx, db.PresenceChangedEvent{
			PlayerID:  id,
			NotifyIDs: friends,
			Online:    false,
		})

		if err != nil {
			log.Println("Error sending disconnect notification:", err)
		}
	}()

	return nil
}
