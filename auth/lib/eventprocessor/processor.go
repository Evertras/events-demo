package eventprocessor

import (
	"bytes"
	"context"
	"log"

	"github.com/pkg/errors"

	"github.com/Evertras/events-demo/auth/lib/authdb"
	"github.com/Evertras/events-demo/auth/lib/stream"
	"github.com/Evertras/events-demo/auth/lib/stream/authevents"
)

type Processor interface {
	Run(ctx context.Context) error
}

type processor struct {
	db           authdb.Db
	streamReader stream.Reader
}

func New(db authdb.Db, streamReader stream.Reader) Processor {
	return &processor{
		db:           db,
		streamReader: streamReader,
	}
}

func (p *processor) Run(ctx context.Context) error {
	p.streamReader.RegisterHandler(stream.EventIDUserRegistered, func(data []byte) error {
		ev, err := authevents.DeserializeUserRegistered(bytes.NewReader(data))

		if err != nil {
			return err
		}

		if ev == nil {
			return errors.New("nil deserialized registration event")
		}

		log.Println("Got registered event for", ev.Email)

		err = p.db.CreateUser(authdb.UserEntry{
			ID:                   ev.ID,
			Email:                ev.Email,
			PasswordHashWithSalt: ev.PasswordHash,
		})

		return err
	})

	// We don't actually do anything synchronously, so just wait for context
	<-ctx.Done()

	return nil
}
