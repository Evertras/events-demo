package authdb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Db interface {
	Connect() error
	Ping() error
	MigrateToLatest() error

	RegisterUser(username string, password string) error
	ValidateUser(username string, password string) (bool, error)
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

// MigrateToLatest applies all migration scripts in order to make
// sure the database schema is up to date.
func (d *db) MigrateToLatest() error {
	index, err := d.getLatestMigrationIndex()

	if err != nil {
		return err
	}

	if index == len(migrations) {
		log.Println("Database is up to date")
		return nil
	}

	if index < 0 || index > len(migrations) {
		return errors.New(fmt.Sprintf("index %d out of range", index))
	}

	tx, err := d.db.Begin()

	if err != nil {
		return err
	}

	log.Println("Performing migrations")
	for i := index; i < len(migrations); i++ {
		log.Println("Performing migration", i)
		err = migrations[i](tx)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Exec(`
UPDATE migration SET last=$1
`, len(migrations))

	return tx.Commit()
}

func (d *db) getLatestMigrationIndex() (int, error) {
	row := d.db.QueryRow(`SELECT last FROM migration`)

	var rawIndex int64

	err := row.Scan(&rawIndex)

	if err != nil && err.Error() == "pq: relation \"migration\" does not exist" {
		rawIndex = 0
		err = nil
	}

	return int(rawIndex), err
}

func (d *db) RegisterUser(username string, password string) error {
	return nil
}

func (d *db) ValidateUser(username string, password string) (bool, error) {
	return true, nil
}
