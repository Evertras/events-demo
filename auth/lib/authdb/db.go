package authdb

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type UserID string

type Db interface {
	Connect() error
	Ping() error
	MigrateToLatest() error

	RegisterUser(email string, password string) (UserID, error)
	ValidateUser(email string, password string) (bool, error)
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
