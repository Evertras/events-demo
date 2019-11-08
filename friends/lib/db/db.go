package db

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/neo4j"
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
	CreatePlayer(ctx context.Context, userID string) error

	// SendInvite sends a new invitation to a target player
	SendInvite(ctx context.Context, t time.Time, fromID string, toID string) error

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

	val, got := result.Record().Get("v.value")

	if !got {
		return "", errors.New("did not return a value")
	}

	return val.(string), nil
}

func (d *db) Close() error {
	d.session.Close()
	d.driver.Close()

	return nil
}

func (d *db) CreatePlayer(ctx context.Context, userID string) error {
	_, err := d.session.Run(
		`MERGE (:Player { id: $userID })`,
		map[string]interface{}{"userID": userID},
	)

	if err != nil {
		return errors.Wrap(err, "failed to write to db")
	}

	return nil
}

func (d *db) SendInvite(ctx context.Context, t time.Time, fromID string, toID string) error {
	_, err := d.session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
MATCH (fromPlayer:Player {id: $fromID})
MATCH (toPlayer:Player {id: $toID})
MERGE (fromPlayer)-[:INVITED { unixSeconds: $unixSeconds }]->(toPlayer)
`,
			map[string]interface{}{"fromID": fromID, "toID": toID, "unixSeconds": t.Unix()},
		)

		if err != nil {
			return nil, err
		}

		return nil, result.Err()
	})

	if err != nil {
		return errors.Wrap(err, "failed to write to db")
	}

	return nil
}

func (d *db) GetPendingInvites(ctx context.Context, id string) ([]string, error) {
	return nil, errors.New("not implemented")
}
