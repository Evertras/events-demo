package friendlist

import (
	"context"
	"testing"

	"github.com/Evertras/events-demo/presence/lib/db/mockdb"
)

func makeMockDb() *mockdb.MockDb {
	return mockdb.New()
}

func TestPlayerHasFriendsList(t *testing.T) {
	db := makeMockDb()
	playerId := "abcdefplayer"
	playerFriends := []string{
		"aslkjgsdf",
		"89134109JDFAdbasdf",
	}

	connectedFriend := playerFriends[0]

	f := New(db)

	db.SetFriendList(playerId, playerFriends)
	db.Connect(connectedFriend)

	returnedFriends, err := f.GetFriendStatus(context.Background(), playerId)

	if err != nil {
		t.Fatal(err)
	}

	if len(returnedFriends) != len(playerFriends) {
		t.Fatalf("Expected %d friends but got %d", len(playerFriends), len(returnedFriends))
	}

	for _, friend := range returnedFriends {
		shouldBeOnline := friend.Id == connectedFriend

		if friend.Online != shouldBeOnline {
			t.Fatalf("Expected connected: %t, but was connected: %t", shouldBeOnline, friend.Online)
		}
	}
}
