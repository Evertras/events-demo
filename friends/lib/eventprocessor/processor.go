package eventprocessor

import (
	"github.com/Evertras/events-demo/friends/lib/db"
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
	return nil
}
