package friendlist

import (
	"context"
	"testing"
	"time"

	"github.com/Evertras/events-demo/presence/lib/db/mockdb"
)

func makeMockDb() *mockdb.MockDb {
	return mockdb.New()
}

func TestPlayerHasFriendsList(t *testing.T) {
	db := makeMockDb()
	defer db.Close()
	playerID := "abcdefplayer"
	playerFriends := []string{
		"aslkjgsdf",
		"89134109JDFAdbasdf",
	}

	connectedFriend := playerFriends[0]

	f := New(db)

	db.SetFriendList(playerID, playerFriends)
	db.Connect(connectedFriend)

	returnedFriends, err := f.GetFriendStatus(context.Background(), playerID)

	if err != nil {
		t.Fatal(err)
	}

	if len(returnedFriends) != len(playerFriends) {
		t.Fatalf("Expected %d friends but got %d", len(playerFriends), len(returnedFriends))
	}

	for _, friend := range returnedFriends {
		shouldBeOnline := friend.ID == connectedFriend

		if friend.Online != shouldBeOnline {
			t.Fatalf("Expected connected: %t, but was connected: %t", shouldBeOnline, friend.Online)
		}
	}
}

func TestNotificationsCloseGracefully(t *testing.T) {
	db := makeMockDb()

	done := make(chan bool)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	f := New(db)

	go func() {
		err := f.ListenForChanges(ctx)

		if err != nil {
			t.Error(err)
		}

		done <- true
	}()

	time.Sleep(time.Millisecond * 10)

	db.Close()

	select {
	case <-done:

	case <-time.After(time.Millisecond * 100):
		t.Error("Did not cleanly exit ListenForChanges")
	}
}

func TestNotificationsSentWhenPresenceChanges(t *testing.T) {
	db := makeMockDb()
	defer db.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playerID := "A"
	playerFriends := []string{
		"B",
	}

	connectedFriend := "B"

	f := New(db)

	go func() {
		if err := f.ListenForChanges(ctx); err != nil {
			t.Error(err)
		}
	}()

	notifications, err := f.Subscribe(ctx, playerID)

	if err != nil {
		t.Fatal("Failed to get notification channel", err)
	}

	db.SetFriendList(playerID, playerFriends)
	db.Connect(connectedFriend)

	select {
	case n := <-notifications:
		if n.ID != connectedFriend {
			t.Errorf("Connect: Expected id %q but found %q", connectedFriend, n.ID)
		}

		if !n.Online {
			t.Errorf("Connect: Expected to be online, but was offline")
		}

	case <-time.After(time.Millisecond * 100):
		t.Error("Connect: No notification found")
	}

	db.Disconnect(connectedFriend)

	select {
	case n := <-notifications:
		if n.ID != connectedFriend {
			t.Errorf("Disconnect: Expected id %q but found %q", connectedFriend, n.ID)
		}

		if n.Online {
			t.Errorf("Disconnect: Expected to be offline, but was online")
		}

	case <-time.After(time.Millisecond * 100):
		t.Error("Disconnect: No notification found")
	}
}
