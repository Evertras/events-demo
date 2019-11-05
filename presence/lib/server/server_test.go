package server

import (
	"context"
	"testing"
	"time"

	"github.com/Evertras/events-demo/presence/lib/connection/mockconnection"
	"github.com/Evertras/events-demo/presence/lib/db/mockdb"
	"github.com/Evertras/events-demo/presence/lib/friendlist"
	"github.com/Evertras/events-demo/presence/lib/server/mocklistener"
)

func TestPlayerNotifiedWhenFriendConnects(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Dependencies
	listener := mocklistener.New()
	db := mockdb.New()
	fl := friendlist.New(db)

	// The actual server
	s := New(listener, db, fl)

	go fl.ListenForChanges(ctx)
	go s.Run(ctx)

	idA := "PlayerA"
	idB := "PlayerB"

	db.MakeFriends(idA, idB)

	connA := mockconnection.New(idA)
	connB := mockconnection.New(idB)

	listener.AddConnection(connA)

	if len(connA.GetReceivedFriendStatusNotifications()) != 0 {
		t.Fatal("Somehow already received notifications, shouldn't have")
	}

	listener.AddConnection(connB)

	time.Sleep(time.Millisecond * 1000)

	rcv := connA.GetReceivedFriendStatusNotifications()

	if len(rcv) != 1 {
		t.Fatalf("Expected 1 received notification but got %d", len(rcv))
	}
}
