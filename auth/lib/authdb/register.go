package authdb

import (
	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

func (d *db) RegisterUser(email string, password string) (UserID, error) {
	tx, err := d.db.Begin()

	if err != nil {
		return "", errors.Wrap(err, "could not create transaction")
	}

	id := UserID(uuid.New().String())

	// bcrypt package takes care of salting for us
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", errors.Wrap(err, "failed to generate hash")
	}

	_, err = tx.Exec(`
INSERT INTO users (id, email, hash)
VALUES ($1, $2, $3)
`,
		id, email, hash)

	if err != nil {
		return "", errors.Wrap(err, "insert query failed")
	}

	err = tx.Commit()

	if err != nil {
		return "", errors.Wrap(err, "failed to commit transaction")
	}

	return id, nil
}
