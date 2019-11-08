package eventprocessor

import (
	"bytes"
	"context"

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

			err = p.db.CreatePlayer(ctx, ev.ID)

			return err
		},
	)

	return nil
}
