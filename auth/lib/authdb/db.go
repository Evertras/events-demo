package authdb

import (
	"database/sql"
	"github.com/pkg/errors"
	"fmt"

	_ "github.com/lib/pq"
)

// UserEntry is a single row in the database
type UserEntry struct {
	// ID is some uniquely generated identifier attached to the user
	ID           string

	// Email is the user's email address
	Email        string

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

	// MigrateToLatest checks the current database version
	// and performs any necessary migration scripts
	MigrateToLatest() error

	// CreateUser creates a user in the database
	CreateUser(entry UserEntry) error

	// GetUserByEmail returns the user's entry or nil if not found.
	// Returns an error if something unexpected goes wrong.
	GetUserByEmail(email string) (*UserEntry, error)
}

type db struct {
	db   *sql.DB
	opts ConnectionOptions
}

type ConnectionOptions struct {
	User     string
	Password string
	Address  string
}

func New(opts ConnectionOptions) Db {
	return &db{
		db:   nil,
		opts: opts,
	}
}

func (d *db) Connect() error {
	var err error

	connStr := fmt.Sprintf("postgres://%s:%s@%s/auth?sslmode=disable", d.opts.User, d.opts.Password, d.opts.Address)

	d.db, err = sql.Open("postgres", connStr)

	return err
}

func (d *db) Ping() error {
	if d.db == nil {
		return errors.New("db not connected")
	}

	return d.db.Ping()
}

func (d *db) CreateUser(entry UserEntry) error {
	tx, err := d.db.Begin()

	if err != nil {
		return errors.Wrap(err, "could not create transaction")
	}

	_, err = tx.Exec(`
INSERT INTO users (id, email, hash)
VALUES ($1, $2, $3)
`,
		entry.ID, entry.Email, entry.PasswordHashWithSalt)

	if err != nil {
		return errors.Wrap(err, "insert query failed")
	}

	err = tx.Commit()

	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

func (d *db) GetUserByEmail(email string) (*UserEntry, error) {
	row := d.db.QueryRow(`
SELECT id, email, hash
FROM users
WHERE email = $1
`, email)

	entry := &UserEntry{}

	err := row.Scan(&entry.ID, &entry.Email, &entry.PasswordHashWithSalt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "unexpected error when scanning row")
	}

	return entry, nil
}
