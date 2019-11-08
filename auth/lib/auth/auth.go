package auth

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/Evertras/events-demo/auth/lib/authdb"
	"github.com/Evertras/events-demo/auth/lib/events"
	"github.com/Evertras/events-demo/auth/lib/events/authevents"
)

var ErrUserAlreadyExists = errors.New("user already exists")

// Auth performs auth operations and updates an underlying data store and event stream
type Auth interface {
	// Register creates a UserRegistered event and waits for the user to be added.
	// Returns the generated user ID for the user on success.
	Register(ctx context.Context, email string, password string) (string, error)

	// Validate checks if the id and password are correct.
	//
	// Returns true if they match
	// Returns false if they do not match, but the check itself was made
	// Returns an error if the check could not be made
	ValidateByID(ctx context.Context, id string, password string) (bool, error)

	// GetIDFromEmail gets the canonical user ID from the given email, if it exists
	//
	// Returns the ID on match
	// Returns an empty string if not found
	// Returns an error if something unexpected occurred
	GetIDFromEmail(ctx context.Context, email string) (string, error)
}

type auth struct {
	db          authdb.Db
	eventWriter events.Writer
}

func New(db authdb.Db, eventWriter events.Writer) Auth {
	a := &auth{
		db:          db,
		eventWriter: eventWriter,
	}

	return a
}

func (a *auth) Register(ctx context.Context, email string, password string) (string, error) {
	fullSpan, ctx := opentracing.StartSpanFromContext(ctx, "Register")
	fullSpan.SetTag("component", "logic")
	defer fullSpan.Finish()

	// Note that this is a best-effort sanity check; if two register commands are
	// sent quickly back to back, this will NOT stop multiple events from being
	// created, and that's okay.  We need to cleanly handle multiple register events
	// later on in the process.  It's still good to filter what we can at this point.
	existingID, err := a.db.GetIDByEmail(ctx, email)

	if err != nil {
		return "", errors.Wrap(err, "failed to check for existing user")
	}

	if existingID != "" {
		return "", ErrUserAlreadyExists
	}

	hashSpan := opentracing.StartSpan("Hash password", opentracing.ChildOf(fullSpan.Context()))
	// bcrypt package takes care of salting for us
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashSpan.Finish()

	if err != nil {
		return "", errors.Wrap(err, "unable to hash password")
	}

	id := uuid.New().String()

	done := make(chan bool)
	errs := make(chan error)

	go func() {
		err := a.db.WaitForCreateUser(ctx, email)

		if err != nil {
			errs <- err
		} else {
			done <- true
		}
	}()

	ev := authevents.NewUserRegistered()

	ev.ID = id
	ev.Email = email
	ev.PasswordHash = string(hash)
	ev.TimeUnixMs = time.Now().Unix()

	err = a.eventWriter.PostRegisteredEvent(ctx, ev)

	if err != nil {
		return "", err
	}

	select {
	case <-done:
		break

	case e := <-errs:
		fullSpan.SetTag("error", true)
		fullSpan.SetTag("error.object", e)
		return "", errors.Wrap(err, "failed to find registered event")

	case <-ctx.Done():
		err := errors.New("context ended")
		fullSpan.SetTag("error", true)
		fullSpan.SetTag("error.object", err)
		return "", err
	}

	return id, nil
}

func (a *auth) ValidateByID(ctx context.Context, id string, password string) (bool, error) {
	fullSpan, ctx := opentracing.StartSpanFromContext(ctx, "Validate")
	defer fullSpan.Finish()

	hashSpan := opentracing.StartSpan("Hash password", opentracing.ChildOf(fullSpan.Context()))
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashSpan.Finish()

	if err != nil {
		err = errors.Wrap(err, "failed to hash supplied password")
		fullSpan.SetTag("error", true)
		fullSpan.SetTag("error.object", err)
		return false, err
	}

	entry, err := a.db.GetUserByID(ctx, id)

	if err != nil {
		err = errors.Wrap(err, "unexpected error while finding user")
		fullSpan.SetTag("error", true)
		fullSpan.SetTag("error.object", err)
		return false, err
	}

	if entry == nil {
		return false, nil
	}

	compareSpan := opentracing.StartSpan("Compare hash and password", opentracing.ChildOf(fullSpan.Context()))
	// bcrypt handles the salt in the encoded value for us
	valid := bcrypt.CompareHashAndPassword([]byte(entry.PasswordHashWithSalt), hashed) != nil
	compareSpan.Finish()

	return valid, nil
}

func (a *auth) GetIDFromEmail(ctx context.Context, email string) (string, error) {
	id, err := a.db.GetIDByEmail(ctx, email)

	if err != nil {
		return "", errors.Wrap(err, "failed to find ID")
	}

	return id, nil
}
