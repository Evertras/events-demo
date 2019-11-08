package db

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type Db interface {
	// Connect connects to the database
	Connect(ctx context.Context) error

	// Close closes all open connections or other resources
	Close() error

	// GetSharedValue atomically retrieves a value or sets it to the provided
	// default if it is not already set
	GetSharedValue(ctx context.Context, key string, ifNotSet string) (string, error)

	// CreatePlayer creates a player
	CreatePlayer(ctx context.Context, userID string, email string) error

	// SendInviteByID sends a new invitation to a target player by target's ID
	SendInviteByID(ctx context.Context, t time.Time, fromID string, toID string) error

	// SendInviteByEmail sends a new invitation to a target player by target's email
	SendInviteByEmail(ctx context.Context, t time.Time, fromID string, toEmail string) error

	// GetPendingInvites gets all pending invites for a player, returning
	// the user IDs that the invites were sent from
	GetPendingInvites(ctx context.Context, id string) ([]string, error)
}

type db struct {
	addr string

	driver  neo4j.Driver
	session neo4j.Session
}

func New(addr string) Db {
	return &db{
		addr: addr,
	}
}

func (d *db) Connect(ctx context.Context) error {
	driver, err := neo4j.NewDriver(d.addr, neo4j.NoAuth())

	if err != nil {
		return errors.Wrap(err, "failed to create driver")
	}

	session, err := driver.Session(neo4j.AccessModeWrite)

	if err != nil {
		return errors.Wrap(err, "failed to create session")
	}

	d.driver = driver
	d.session = session

	return nil
}

func (d *db) GetSharedValue(ctx context.Context, key string, ifNotSet string) (string, error) {
	_, err := d.session.Run(
		`CREATE CONSTRAINT ON (v:SharedValue) ASSERT v.key IS UNIQUE`,
		nil,
	)

	if err != nil {
		return "", errors.Wrap(err, "failed to create constraint")
	}

	result, err := d.session.Run(
		`
MERGE (v:SharedValue { key: $key })
ON CREATE SET v.value = $value
RETURN v.value
		`,
		map[string]interface{}{"key": key, "value": ifNotSet},
	)

	if err != nil {
		return "", errors.Wrap(err, "failed to run query")
	}

	if result.Err() != nil {
		return "", errors.Wrap(result.Err(), "result failed")
	}

	for result.Next() {
		val, got := result.Record().Get("v.value")

		if !got {
			return "", errors.New("did not return a value in record")
		}

		return val.(string), nil
	}

	return "", errors.New("did not find any records")
}

func (d *db) Close() error {
	d.session.Close()
	d.driver.Close()

	return nil
}

func (d *db) CreatePlayer(ctx context.Context, userID string, email string) error {
	span, ctx := startSpan(ctx, "Create Player")
	defer span.Finish()

	_, err := d.session.Run(
		`MERGE (:Player { playerID: $userID, email: $email })`,
		map[string]interface{}{
			"userID": userID,
			"email":  email,
		},
	)

	if err != nil {
		return errors.Wrap(err, "failed to write to db")
	}

	return nil
}

func (d *db) SendInviteByID(ctx context.Context, t time.Time, fromID string, toID string) error {
	span, ctx := startSpan(ctx, "Record Invite By ID")
	defer span.Finish()

	span.SetTag("from", fromID)
	span.SetTag("target", toID)

	result, err := d.session.Run(
		`
MATCH (fromPlayer:Player { playerID: $fromID })
MATCH (toPlayer:Player { playerID: $toID })
MERGE (fromPlayer)-[i:INVITED]->(toPlayer)
ON CREATE SET i.time = $time
ON MATCH SET i.time = $time
`,
		map[string]interface{}{"fromID": fromID, "toID": toID, "time": t},
	)

	if err != nil {
		return errors.Wrap(err, "failed to write to db")
	}

	summary, err := result.Summary()

	if err != nil {
		return errors.Wrap(err, "failed to get result summary")
	}

	if summary.Counters().RelationshipsCreated() != 1 {
		return errors.New("did not create relationship")
	}

	return nil
}

func (d *db) SendInviteByEmail(ctx context.Context, t time.Time, fromID string, toEmail string) error {
	span, ctx := startSpan(ctx, "Record Invite By Email")
	defer span.Finish()

	span.SetTag("from", fromID)
	span.SetTag("target", toEmail)

	result, err := d.session.Run(
		`
MATCH (fromPlayer:Player { playerID: $fromID })
MATCH (toPlayer:Player { email: $toEmail })
MERGE (fromPlayer)-[i:INVITED]->(toPlayer)
ON CREATE SET i.time = $time
ON MATCH SET i.time = $time
`,
		map[string]interface{}{"fromID": fromID, "toEmail": toEmail, "time": t},
	)

	if err != nil {
		return errors.Wrap(err, "failed to write to db")
	}

	summary, err := result.Summary()

	if err != nil {
		return errors.Wrap(err, "failed to get result summary")
	}

	if summary.Counters().RelationshipsCreated() != 1 {
		return errors.New("did not create relationship")
	}

	return nil
}

func (d *db) GetPendingInvites(ctx context.Context, id string) ([]string, error) {
	return nil, errors.New("not implemented")
}

func startSpan(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName)

	span.SetTag("db.type", "neo4j")
	span.SetTag("span.kind", "client")
	span.SetTag("component", "db")

	return span, ctx
}
