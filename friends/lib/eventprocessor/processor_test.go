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

func genMockInviteEvent(fromID string, toID string, toEmail string) mockstream.MockReceivedEvent {
	ev := friendevents.NewInviteSent()

	ev.FromID = fromID
	ev.ToID = toID
	ev.ToEmail = toEmail

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
	sourceID := "sourceid"
	targetID := "anotherPlayer"
	targetEmail := "another@somewhere.com"

	cases := []struct {
		name    string
		toID    string
		toEmail string
	}{
		{
			name:    "ID",
			toID:    targetID,
			toEmail: "",
		},
		{
			name:    "Email",
			toID:    "",
			toEmail: targetEmail,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			db := mockdb.New()
			p := New(db)

			db.CreatePlayer(ctx, sourceID, "doesnt@matter.com")
			db.CreatePlayer(ctx, targetID, targetEmail)

			streamReader := mockstream.NewReader()

			go streamReader.Listen(ctx)

			p.RegisterHandlers(streamReader)

			streamReader.Receive <- genMockInviteEvent(sourceID, c.toID, c.toEmail)

			time.Sleep(time.Millisecond * 100)

			if !db.MockInviteExists(sourceID, targetID, targetEmail) {
				t.Error("Did not find invite in DB after invite event was sent")
			}
		})
	}
}
