package authdb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// This file contains some simple migration 'scripts' to perform schema
// migrations to the auth DB safely.  You probably want to do this in
// a more robust way in a real app, but this is simple enough for now.
type migrationFunc = func(tx *sql.Tx) error

func migration_000_MigrationTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
CREATE TABLE migration (
	last integer NOT NULL
);

INSERT INTO migration(last) VALUES (0);
`)

	return err
}

func migration_001_UserTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
CREATE TABLE users (
	id VARCHAR(128) NOT NULL PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	hash TEXT NOT NULL
);
`)

	return err
}

var migrations = []migrationFunc{
	migration_000_MigrationTable,
	migration_001_UserTable,
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
