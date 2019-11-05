package mockdb

import (
	"context"
	"testing"
	"time"
)

// NOTE: Yes, this is a test for testing.  Who watches the watchmen?
//       Unfortunately a little testing of the tests is necessary to
//       have confidence in the "real" tests.

func TestConnectSendsNotifications(t *testing.T) {
	ctx := context.Background()
	m := New()

	playerID := "some-id"
	notifyList := []string{
		"A",
		"B",
	}

	m.SetFriendList(ctx, playerID, notifyList)

	gotNotification := make(chan bool, 1)

	go func() {
		select {
		case n := <-m.notifications:

			if n.PlayerID != playerID {
				t.Errorf("Notification: Expected %q but got %q", playerID, n.PlayerID)
			}

			if len(n.NotifyIDs) != len(notifyList) {
				t.Errorf("Notification: Expected %d notify IDs but got %d", len(notifyList), len(n.NotifyIDs))
			}

			gotNotification <- true

		case <-time.After(time.Millisecond * 100):
			t.Error("Notification: Never got notification")
		}
	}()

	m.Connect(playerID)

	select {
	case <-gotNotification:

	case <-time.After(time.Millisecond * 100):
		t.Error("Waiting: Did not get notification")
	}
}
