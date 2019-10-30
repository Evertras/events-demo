package eventprocessor

import (
	"bytes"
	"context"

	"github.com/pkg/errors"

	"github.com/Evertras/events-demo/auth/lib/authdb"
	"github.com/Evertras/events-demo/auth/lib/stream"
	"github.com/Evertras/events-demo/auth/lib/stream/authevents"
)

type Processor interface {
	RegisterHandlers(streamReader stream.Reader) error
}

type processor struct {
	db authdb.Db
}

func New(db authdb.Db) Processor {
	return &processor{
		db: db,
	}
}

func (p *processor) RegisterHandlers(streamReader stream.Reader) error {
	return streamReader.RegisterHandler(stream.EventIDUserRegistered, func(ctxInner context.Context, data []byte) error {
		ev, err := authevents.DeserializeUserRegistered(bytes.NewReader(data))

		if err != nil {
			return err
		}

		if ev == nil {
			return errors.New("nil deserialized registration event")
		}

		err = p.db.CreateUser(ctxInner, authdb.UserEntry{
			ID:                   ev.ID,
			Email:                ev.Email,
			PasswordHashWithSalt: ev.PasswordHash,
		})

		return err
	})
}
