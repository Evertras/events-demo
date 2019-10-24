package auth

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/Evertras/events-demo/auth/lib/authdb"
	"github.com/Evertras/events-demo/auth/lib/eventstream"
	"github.com/Evertras/events-demo/auth/lib/eventstream/events"
)

var ErrUserAlreadyExists = errors.New("user already exists")

// RegistrationMeta contains metadata about the user that auth
// doesn't necessarily care about.  However, we still want to be
// explicit about the schema for type safety.
type RegistrationMeta struct {
	Username string
}

// Auth performs auth operations and updates an underlying data store and event stream
type Auth interface {
	// Register adds a new user and creates a registration event
	Register(email string, password string, details RegistrationMeta) error

	// Validate checks if the email and password are correct.
	//
	// Returns true if they match
	// Returns false if they do not match, but the check itself was made
	// Returns an error if the check could not be made
	Validate(email string, password string) (bool, error)
}

type auth struct {
	db authdb.Db
	es eventstream.EventStream
}

func New(db authdb.Db, es eventstream.EventStream) Auth {
	return &auth{
		db: db,
		es: es,
	}
}

func (a *auth) Register(email string, password string, details RegistrationMeta) error {
	exists, err := a.db.GetUserByEmail(email)

	if err != nil {
		return errors.Wrap(err, "failed to check for existing user")
	}

	if exists != nil {
		return ErrUserAlreadyExists
	}

	// bcrypt package takes care of salting for us
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return errors.Wrap(err, "unable to hash password")
	}

	entry := authdb.UserEntry{
		ID:                   uuid.New().String(),
		Email:                email,
		PasswordHashWithSalt: string(hash),
	}

	err = a.db.CreateUser(entry)

	if err != nil {
		return errors.Wrap(err, "unable to create user")
	}

	ev := events.NewUserRegistered()

	ev.ID = entry.ID
	ev.Email = entry.Email
	ev.Username = details.Username

	a.es.PostRegisteredEvent(ev)

	return nil
}

func (a *auth) Validate(email string, password string) (bool, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return false, errors.Wrap(err, "failed to hash supplied password")
	}

	entry, err := a.db.GetUserByEmail(email)

	if err != nil {
		return false, errors.Wrap(err, "unexpected error while finding user")
	}

	if entry == nil {
		return false, nil
	}

	// bcrypt handles the salt in the encoded value for us
	return bcrypt.CompareHashAndPassword([]byte(entry.PasswordHashWithSalt), hashed) != nil, nil
}
