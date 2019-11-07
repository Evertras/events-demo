package auth

import (
	"context"
	"testing"
	"time"

	"github.com/Evertras/events-demo/auth/lib/authdb/mockauthdb"
	"github.com/Evertras/events-demo/auth/lib/events"
	"github.com/Evertras/events-demo/auth/lib/events/authevents"
	mockstream "github.com/Evertras/events-demo/shared/stream/mock"
)

func TestRegisterUserSendsCompleteRegisterEventWithHashedID(t *testing.T) {
	db := mockauthdb.New()
	inner := mockstream.NewWriter()

	writer := events.NewWriter(inner)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	a := New(db, writer)

	email := "test@testing.com"
	password := "sekrit"

	id, err := a.Register(ctx, email, password)

	if err != nil {
		t.Fatal("Failed to register", err)
	}

	if len(inner.Sent) != 1 {
		t.Fatalf("Expected 1 event sent, but found %d", len(inner.Sent))
	}

	ev := inner.Sent[0].Payload.(*authevents.UserRegistered)

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
