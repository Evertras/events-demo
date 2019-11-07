package eventprocessor

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/Evertras/events-demo/auth/lib/authdb/mockauthdb"
	"github.com/Evertras/events-demo/auth/lib/events"
	"github.com/Evertras/events-demo/auth/lib/events/authevents"
	mockstream "github.com/Evertras/events-demo/shared/stream/mock"
)

func TestRegistrationCreatesUserBasedOnRegistrationEvent(t *testing.T) {
	db := mockauthdb.New()
	reader := mockstream.NewReader()

	r := New(db)

	err := r.RegisterHandlers(reader)

	if err != nil {
		t.Fatal("Failed to register handlers:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	go reader.Listen(ctx)

	receivedEvent := authevents.NewUserRegistered()

	id := "abcdefid"
	email := "test803451234@testgajgklawe.com"
	passwordHash := "hasdhasdhasdfhashgkjiw"

	receivedEvent.Email = email
	receivedEvent.PasswordHash = passwordHash
	receivedEvent.ID = id

	var buf bytes.Buffer

	receivedEvent.Serialize(&buf)

	reader.Receive <- mockstream.MockReceivedEvent{
		ID:   events.EventIDUserRegistered,
		Data: buf.Bytes(),
	}

	time.Sleep(time.Millisecond * 10)

	if len(db.EntriesByID) == 0 {
		t.Fatal("No users added at all")
	}

	user, err := db.GetUserByEmail(context.Background(), email)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("User not found")
	}

	if user.ID != id {
		t.Errorf("Expected user ID of %q but found %q", id, user.ID)
	}

	if user.Email != email {
		t.Errorf("Expected user email of %q but found %q", email, user.Email)
	}

	if user.PasswordHashWithSalt != passwordHash {
		t.Errorf("Expected password hash of %q but found %q", passwordHash, user.PasswordHashWithSalt)
	}
}
