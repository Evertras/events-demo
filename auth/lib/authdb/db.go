package authdb

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis"
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
	Connect() error

	// Ping pings the database to check connectivity
	Ping() error

	// CreateUser creates a user in the database
	CreateUser(entry UserEntry) error

	// GetUserByEmail returns the user's entry or nil if not found.
	// Returns an error if something unexpected goes wrong.
	GetUserByEmail(email string) (*UserEntry, error)

	// GetSharedID returns a stored string ID, or, if it does not exist,
	// creates a random ID, stores it in the DB, and returns it
	GetSharedID(key string) (string, error)

	// WaitForCreateUser will wait for the user to be created
	// in the database before returning, or an error if the user
	// was not seen as added before the context cancels
	WaitForCreateUser(ctx context.Context, email string) error
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

func (d *db) Connect() error {
	d.db = redis.NewClient(&redis.Options{
		Addr:     d.opts.Address,
		Password: "",
		DB:       0,
	})

	return d.Ping()
}

func (d *db) Ping() error {
	if d.db == nil {
		return errors.New("db not connected")
	}

	return d.db.Ping().Err()
}

func (d *db) CreateUser(entry UserEntry) error {
	if entry.Email == "" {
		return errors.New("must supply email")
	}

	val, err := json.Marshal(entry)

	if err != nil {
		return errors.Wrap(err, "failed to marshal to JSON")
	}

	err = d.db.Set(credsKey(entry.Email), val, 0).Err()

	return err
}

func (d *db) WaitForCreateUser(ctx context.Context, email string) error {
	ps := d.db.PSubscribe("__keyspace@*__:creds:" + email)

	_, err := ps.Receive()

	if err != nil {
		return err
	}

	defer ps.Close()

	select {
	case <-ps.Channel():
		return nil

	case <-ctx.Done():
		return errors.New("context finished before message received")
	}
}

func (d *db) GetUserByEmail(email string) (*UserEntry, error) {
	entry := &UserEntry{}

	raw, err := d.db.Get(credsKey(email)).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to get key")
	}

	err = json.Unmarshal([]byte(raw), entry)

	if err != nil {
		return nil, errors.Wrap(err, "found key but could not unmarshal json")
	}

	return entry, nil
}

func credsKey(email string) string {
	return "creds:" + email
}

func (d *db) GetSharedID(key string) (string, error) {
	randomID := uuid.New().String()

	locker := redislock.New(d.db)

	lockKey := key + ".lock"

	lock, err := locker.Obtain(
		lockKey,
		50*time.Millisecond,
		nil,
	)

	if err != nil {
		return "", err
	}

	defer lock.Release()

	d.db.SetNX(key, randomID, 0)

	actualID, err := d.db.Get(key).Result()

	return actualID, err
}
