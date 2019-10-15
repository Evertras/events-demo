package authdb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Db struct {
	db   *sql.DB
	opts ConnectionOptions
}

type ConnectionOptions struct {
	User     string
	Password string
	Address  string
}

func New(opts ConnectionOptions) *Db {
	return &Db{
		db:   nil,
		opts: opts,
	}
}

func (db *Db) Connect() error {
	var err error

	connStr := fmt.Sprintf("postgres://%s:%s@%s/auth?sslmode=disable", db.opts.User, db.opts.Password, db.opts.Address)

	db.db, err = sql.Open("postgres", connStr)

	return err
}

func (db *Db) Ping() error {
	if db.db == nil {
		return errors.New("db not connected")
	}

	return db.db.Ping()
}

// MigrateToLatest applies all migration scripts in order to make
// sure the database schema is up to date.
func (db *Db) MigrateToLatest() error {
	index, err := db.getLatestMigrationIndex()

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

	tx, err := db.db.Begin()

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

func (db *Db) getLatestMigrationIndex() (int, error) {
	row := db.db.QueryRow(`SELECT last FROM migration`)

	var rawIndex int64

	err := row.Scan(&rawIndex)

	if err != nil && err.Error() == "pq: relation \"migration\" does not exist" {
		rawIndex = 0
		err = nil
	}

	return int(rawIndex), err
}
