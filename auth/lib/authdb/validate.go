package authdb

import (
	"database/sql"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (d *db) ValidateUser(email string, password string) (bool, error) {
	row := d.db.QueryRow(`
SELECT hash
FROM users
WHERE email = $1
`, email)

	var stored string

	err := row.Scan(&stored)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return false, errors.Wrap(err, "failed to hash supplied password")
	}

	return bcrypt.CompareHashAndPassword([]byte(stored), hashed) != nil, nil
}
