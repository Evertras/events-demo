package eventprocessor

import (
	"bytes"
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/Evertras/events-demo/friends/lib/db"
	"github.com/Evertras/events-demo/friends/lib/events"
	"github.com/Evertras/events-demo/friends/lib/events/friendevents"
	"github.com/Evertras/events-demo/shared/stream"
)

type Processor interface {
	RegisterHandlers(streamReader stream.Reader) error
}

type processor struct {
	db db.Db
}

func New(db db.Db) Processor {
	return &processor{
		db: db,
	}
}

func (p *processor) RegisterHandlers(streamReader stream.Reader) error {
	streamReader.RegisterHandler(
		events.EventIDUserRegistered,
		func(ctx context.Context, data []byte) error {
			ev, err := friendevents.DeserializeUserRegistered(bytes.NewReader(data))

			if err != nil {
				return err
			}

			if ev == nil {
				return errors.New("nil deserialized registration event")
			}

			err = p.db.CreatePlayer(ctx, ev.ID, ev.Email)

			return err
		},
	)

	streamReader.RegisterHandler(
		events.EventIDInviteSent,
		func(ctx context.Context, data []byte) error {
			ev, err := friendevents.DeserializeInviteSent(bytes.NewReader(data))

			if err != nil {
				return errors.Wrap(err, "failed to deserialize invite message")
			}

			if ev == nil {
				return errors.New("nil deserialized invite sent event")
			}

			if len(ev.ToID) > 0 {
				err = p.db.SendInviteByID(ctx, time.Unix(ev.TimeUnixMs, 0), ev.FromID, ev.ToID)
			} else if len(ev.ToEmail) > 0 {
				err = p.db.SendInviteByEmail(ctx, time.Unix(ev.TimeUnixMs, 0), ev.FromID, ev.ToEmail)
			} else {
				return errors.New("must supply either ID or email for invitation")
			}

			return err
		},
	)

	return nil
}
