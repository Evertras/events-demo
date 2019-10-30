package mockauthdb

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/Evertras/events-demo/auth/lib/authdb"
)

type MockDb struct {
	Connected bool

	EntriesByEmail map[string]*authdb.UserEntry
	EntriesByID    map[string]*authdb.UserEntry

	CreateUserTimeout time.Duration
}

// Ensure it meets interface
var _ authdb.Db = &MockDb{}

func New() *MockDb {
	return &MockDb{
		Connected:         false,
		EntriesByEmail:    make(map[string]*authdb.UserEntry),
		EntriesByID:       make(map[string]*authdb.UserEntry),
		CreateUserTimeout: time.Millisecond,
	}
}

func (db *MockDb) Connect() error {
	db.Connected = true

	return nil
}

func (db *MockDb) Ping() error {
	return nil
}

func (db *MockDb) CreateUser(entry authdb.UserEntry) error {
	if _, exists := db.EntriesByEmail[entry.Email]; exists {
		return errors.New("email already exists")
	}

	if _, exists := db.EntriesByID[entry.ID]; exists {
		return errors.New("id already exists")
	}

	e := &authdb.UserEntry{
		ID:                   entry.ID,
		Email:                entry.Email,
		PasswordHashWithSalt: entry.PasswordHashWithSalt,
	}

	db.EntriesByEmail[entry.Email] = e
	db.EntriesByID[entry.ID] = e

	return nil
}

func (db *MockDb) GetUserByEmail(email string) (*authdb.UserEntry, error) {
	if entry, ok := db.EntriesByEmail[email]; ok {
		return entry, nil
	}

	return nil, nil
}

func (db *MockDb) WaitForCreateUser(ctx context.Context, email string) error {
	time.Sleep(db.CreateUserTimeout)

	return nil
}

func (db *MockDb) GetSharedValue(key string, ifNotExist string) (string, error) {
	return "someval", nil
}
