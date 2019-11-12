package eventprocessor

import (
	"bytes"
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Evertras/events-demo/friends/lib/events"
	"github.com/Evertras/events-demo/friends/lib/events/friendevents"

	mockdb "github.com/Evertras/events-demo/friends/lib/db/mock"
	mockstream "github.com/Evertras/events-demo/shared/stream/mock"
)

func TestAddsRegisteredUsersToDb(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := mockdb.New()
	p := New(db)

	id := "someplayer"

	streamReader := mockstream.NewReader(mockstream.MockStreamReaderOpts{})

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
	sourceEmail := "source@me.com"
	targetID := "anotherPlayerId"
	targetEmail := "another@somewhere.com"

	cases := []struct {
		name      string
		toID      string
		toEmail   string
		shouldAdd bool
	}{
		{
			name:      "ID",
			toID:      targetID,
			toEmail:   "",
			shouldAdd: true,
		},
		{
			name:      "Email",
			toID:      "",
			toEmail:   targetEmail,
			shouldAdd: true,
		},
		{
			name:      "Self-targeted",
			toID:      "",
			toEmail:   sourceEmail,
			shouldAdd: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			db := mockdb.New()
			p := New(db)

			db.CreatePlayer(ctx, sourceID, sourceEmail)
			db.CreatePlayer(ctx, targetID, targetEmail)

			streamReader := mockstream.NewReader(mockstream.MockStreamReaderOpts{
				Logger: log.New(os.Stdout, "SendInvite-" + c.name + " - ", 0),
			})

			go streamReader.Listen(ctx)

			p.RegisterHandlers(streamReader)

			streamReader.Receive <- genMockInviteEvent(sourceID, c.toID, c.toEmail)

			time.Sleep(time.Millisecond * 100)

			if db.MockInviteExists(sourceID, c.toID, c.toEmail) != c.shouldAdd {
				if c.shouldAdd {
					t.Error("Expected to find invite in DB but did not")
				} else {
					t.Error("Did not expect to find invite in DB but did")
				}
			}
		})
	}
}

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
