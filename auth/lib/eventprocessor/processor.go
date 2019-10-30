package eventprocessor

import (
	"bytes"
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/Evertras/events-demo/auth/lib/authdb"
	"github.com/Evertras/events-demo/auth/lib/stream"
	"github.com/Evertras/events-demo/auth/lib/stream/authevents"
	"github.com/Evertras/events-demo/auth/lib/tracing"
)

type Processor interface {
	RegisterHandlers(streamReader stream.Reader) error
}

type processor struct {
	db           authdb.Db
	tracer       opentracing.Tracer
}

func New(db authdb.Db) (Processor, error) {
	tracer, err := tracing.Init("processor")

	if err != nil {
		return nil, errors.Wrap(err, "failed to init tracer")
	}

	return &processor{
		db:           db,
		tracer:       tracer,
	}, nil
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
