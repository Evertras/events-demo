package authdb

import "database/sql"

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
	name VARCHAR(32) NOT NULL PRIMARY KEY,
	hash VARCHAR(128) NOT NULL,
	salt VARCHAR(128) NOT NULL
);
`)

	return err
}

var migrations = []migrationFunc {
	migration_000_MigrationTable,
	migration_001_UserTable,
}
