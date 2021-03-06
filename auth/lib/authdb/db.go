package authdb

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// UserEntry is a single row in the database
type UserEntry struct {
	// ID is some uniquely generated identifier attached to the user
	ID string

	// Email is the user's email address
	Email string

	// PasswordHash is the hashed/salted password that will be used
	// for comparison/validation
	PasswordHashWithSalt string
}

// Db is a persistent database that stores user information
type Db interface {
	// Connect actually connects to the database
	Connect(ctx context.Context) error

	// Ping pings the database to check connectivity
	Ping(ctx context.Context) error

	// CreateUser creates a user in the database
	CreateUser(ctx context.Context, entry UserEntry) error

	// GetIDByEmail returns the user's ID or empty string if not found.
	// Returns an error if something unexpected goes wrong.
	GetIDByEmail(ctx context.Context, email string) (string, error)

	// GetUserByID returns the user's entry or nil if not found.
	// Returns an error if something unexpected goes wrong.
	GetUserByID(ctx context.Context, id string) (*UserEntry, error)

	// GetSharedValue returns a stored value; if it does not exist,
	// it will atomically store the given value and return that.
	GetSharedValue(ctx context.Context, key string, ifNotExist string) (string, error)

	// WaitForCreateUser will wait for the user to be created
	// in the database before returning, or an error if the user
	// was not seen as added before the context cancels
	WaitForCreateUser(ctx context.Context, id string) error
}

type db struct {
	db   *redis.Client
	opts ConnectionOptions
}

type ConnectionOptions struct {
	Address string
}

func New(opts ConnectionOptions) Db {
	return &db{
		db:   nil,
		opts: opts,
	}
}

func (d *db) Connect(ctx context.Context) error {
	d.db = redis.NewClient(&redis.Options{
		Addr:     d.opts.Address,
		Password: "",
		DB:       0,
	})

	return d.Ping(ctx)
}

func startSpan(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName)

	span.SetTag("db.type", "redis")
	span.SetTag("span.kind", "client")
	span.SetTag("component", "authdb")

	return span, ctx
}

func (d *db) Ping(ctx context.Context) error {
	span, ctx := startSpan(ctx, "Redis Ping")
	span.SetTag("db.statement", "PING")
	defer span.Finish()

	if d.db == nil {
		return errors.New("db not connected")
	}

	return d.db.Ping().Err()
}

func (d *db) CreateUser(ctx context.Context, entry UserEntry) error {
	span, ctx := startSpan(ctx, "Redis CreateUser")
	defer span.Finish()

	if entry.Email == "" {
		return errors.New("must supply email")
	}

	val, err := json.Marshal(entry)

	if err != nil {
		return errors.Wrap(err, "failed to marshal to JSON")
	}

	// TODO: Transactionify this
	err = d.db.Set(keyEmail(entry.Email), entry.ID, 0).Err()

	if err != nil {
		return errors.Wrap(err, "failed to write email key")
	}

	err = d.db.Set(keyID(entry.ID), val, 0).Err()

	if err != nil {
		return errors.Wrap(err, "failed to write ID key")
	}

	return err
}

func (d *db) WaitForCreateUser(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Wait for user creation")
	defer span.Finish()

	ps := d.db.PSubscribe("__keyspace@*__:" + keyID(id))

	defer ps.Close()

	_, err := ps.Receive()

	if err != nil {
		return err
	}

	select {
	case <-ps.Channel():
		return nil

	case <-ctx.Done():
		return errors.New("context finished before message received")
	}
}

func (d *db) GetIDByEmail(ctx context.Context, email string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get ID by email")
	defer span.Finish()

	id, err := d.db.Get(keyEmail(email)).Result()

	if err != nil {
		if err == redis.Nil {
			return "", nil
		}

		err = errors.Wrap(err, "failed to get key")

		span.SetTag("error", true)
		span.SetTag("error.object", err)

		return "", err
	}

	return id, nil
}

func (d *db) GetUserByID(ctx context.Context, id string) (*UserEntry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get user by ID")
	defer span.Finish()

	entry := &UserEntry{}

	raw, err := d.db.Get(keyID(id)).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		err = errors.Wrap(err, "failed to get key")

		span.SetTag("error", true)
		span.SetTag("error.object", err)

		return nil, err
	}

	err = json.Unmarshal([]byte(raw), entry)

	if err != nil {
		err = errors.Wrap(err, "found key but could not unmarshal json")

		span.SetTag("error", true)
		span.SetTag("error.object", err)

		return nil, err
	}

	return entry, nil
}

func keyEmail(email string) string {
	return "email:" + email
}

func keyID(email string) string {
	return "id:" + email
}

func (d *db) GetSharedValue(ctx context.Context, key string, ifNotExist string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get shared value")
	span.SetTag("auth.sharedvalue.key", key)
	defer span.Finish()

	locker := redislock.New(d.db)

	lockKey := key + ".lock"

	lock, err := locker.Obtain(
		lockKey,
		50*time.Millisecond,
		nil,
	)

	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error.object", err)
		return "", err
	}

	defer lock.Release()

	d.db.SetNX(key, ifNotExist, 0)

	actualID, err := d.db.Get(key).Result()

	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error.object", err)
		return "", err
	}

	return actualID, nil
}
