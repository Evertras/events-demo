package eventprocessor

import (
	"bytes"
	"testing"
	"time"

	"github.com/Evertras/events-demo/friends/lib/events"
	"github.com/Evertras/events-demo/friends/lib/events/friendevents"

	mockdb "github.com/Evertras/events-demo/friends/lib/db/mock"
	mockstream "github.com/Evertras/events-demo/shared/stream/mock"
)

func genMockReceiveUserRegisterEvent(id string, email string) mockstream.MockReceivedEvent {
	ev := friendevents.NewUserRegistered()

	ev.ID = id
	ev.Email = email

	var buf bytes.Buffer

	ev.Serialize(&buf)

	return mockstream.MockReceivedEvent{
		ID:   events.EventIDUserRegistered,
		Data: buf.Bytes(),
	}
}

func TestAddsRegisteredUsersToDb(t *testing.T) {
	db := mockdb.New()
	p := New(db)

	id := "someplayer"
	email := "some@player.com"

	streamReader := mockstream.NewReader()

	p.RegisterHandlers(streamReader)

	streamReader.Receive <- genMockReceiveUserRegisterEvent(id, email)

	time.Sleep(time.Millisecond * 100)

	if _, exists := db.Invites[id]; !exists {
		t.Error("Did not find player in DB after registration event was sent")
	}
}
