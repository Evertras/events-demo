package auth

import (
	"context"
	"testing"
	"time"

	"github.com/Evertras/events-demo/auth/lib/authdb/mockauthdb"
	"github.com/Evertras/events-demo/auth/lib/stream/authevents"
	"github.com/Evertras/events-demo/auth/lib/stream/mockstream"
)

func TestRegisterUserSendsCompleteRegisterEventWithHashedID(t *testing.T) {
	db := mockauthdb.New()
	writer := mockstream.NewWriter()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	a, err := New(db, writer)

	if err != nil {
		t.Fatal(err)
	}

	email := "test@testing.com"
	password := "sekrit"

	id, err := a.Register(ctx, email, password)

	if err != nil {
		t.Fatal("Failed to register", err)
	}

	if len(writer.Sent) != 1 {
		t.Fatalf("Expected 1 event sent, but found %d", len(writer.Sent))
	}

	ev := writer.Sent[0].(*authevents.UserRegistered)

	if ev == nil {
		t.Fatal("Expected authevents.UserRegistered but did not cast correctly")
	}

	if id == "" {
		t.Error("Did not generate ID")
	} else if ev.ID != string(id) {
		t.Error("Did not set event ID correctly")
	}

	if ev.Email != email {
		t.Errorf("Did not send correct user email: expected %q but got %q", email, ev.Email)
	}

	if ev.TimeUnixMs == 0 {
		t.Error("Did not set timestamp")
	}
}
