package eventprocessor

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/Evertras/events-demo/friends/lib/events"
	"github.com/Evertras/events-demo/friends/lib/events/friendevents"

	mockdb "github.com/Evertras/events-demo/friends/lib/db/mock"
	mockstream "github.com/Evertras/events-demo/shared/stream/mock"
)

func genMockReceiveUserRegisterEvent(id string) mockstream.MockReceivedEvent {
	ev := friendevents.NewUserRegistered()

	ev.ID = id

	var buf bytes.Buffer

	ev.Serialize(&buf)

	return mockstream.MockReceivedEvent{
		ID:   events.EventIDUserRegistered,
		Data: buf.Bytes(),
	}
}

func genMockInviteEvent(fromID string, toID string) mockstream.MockReceivedEvent {
	ev := friendevents.NewInviteSent()

	ev.FromID = fromID
	ev.ToID = toID

	var buf bytes.Buffer

	ev.Serialize(&buf)

	return mockstream.MockReceivedEvent{
		ID:   events.EventIDInviteSent,
		Data: buf.Bytes(),
	}
}

func TestAddsRegisteredUsersToDb(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := mockdb.New()
	p := New(db)

	id := "someplayer"

	streamReader := mockstream.NewReader()

	go streamReader.Listen(ctx)

	p.RegisterHandlers(streamReader)

	streamReader.Receive <- genMockReceiveUserRegisterEvent(id)

	time.Sleep(time.Millisecond * 100)

	if !db.MockPlayerExists(id) {
		t.Error("Did not find player in DB after registration event was sent")
	}
}

func TestAddsInvitesToDb(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := mockdb.New()
	p := New(db)

	fromID := "somePlayer"
	toID := "anotherPlayer"

	db.CreatePlayer(ctx, fromID, "1")
	db.CreatePlayer(ctx, toID, "2")

	streamReader := mockstream.NewReader()

	go streamReader.Listen(ctx)

	p.RegisterHandlers(streamReader)

	streamReader.Receive <- genMockInviteEvent(fromID, toID)

	time.Sleep(time.Millisecond * 100)

	if !db.MockInviteExists(fromID, toID) {
		t.Error("Did not find invite in DB after invite event was sent")
	}
}
