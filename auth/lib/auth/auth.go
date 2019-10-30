package auth

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/Evertras/events-demo/auth/lib/authdb"
	"github.com/Evertras/events-demo/auth/lib/stream"
	"github.com/Evertras/events-demo/auth/lib/stream/authevents"
	"github.com/Evertras/events-demo/auth/lib/tracing"
)

type UserID string

var ErrUserAlreadyExists = errors.New("user already exists")

// Auth performs auth operations and updates an underlying data store and event stream
type Auth interface {
	// Register creates a UserRegistered event and waits for the user to be added
	Register(ctx context.Context, email string, password string) (UserID, error)

	// Validate checks if the email and password are correct.
	//
	// Returns true if they match
	// Returns false if they do not match, but the check itself was made
	// Returns an error if the check could not be made
	Validate(ctx context.Context, email string, password string) (bool, error)
}

type auth struct {
	db           authdb.Db
	streamWriter stream.Writer
	tracer       opentracing.Tracer
}

func New(db authdb.Db, streamWriter stream.Writer) (Auth, error) {
	// TODO: "logic" is a code smell imo, revisit later
	tracer, err := tracing.Init("logic")

	if err != nil {
		return nil, errors.Wrap(err, "failed to init tracer")
	}

	a := &auth{
		db:           db,
		streamWriter: streamWriter,
		tracer:       tracer,
	}

	return a, nil
}

func (a *auth) Register(ctx context.Context, email string, password string) (UserID, error) {
	fullSpan, ctx := opentracing.StartSpanFromContextWithTracer(ctx, a.tracer, "Register")
	defer fullSpan.Finish()

	// Note that this is a best-effort sanity check; if two register commands are
	// sent quickly back to back, this will NOT stop multiple events from being
	// created, and that's okay.  We need to cleanly handle multiple register events
	// later on in the process.  It's still good to filter what we can at this point.
	exists, err := a.db.GetUserByEmail(ctx, email)

	if err != nil {
		return "", errors.Wrap(err, "failed to check for existing user")
	}

	if exists != nil {
		return "", ErrUserAlreadyExists
	}

	hashSpan := a.tracer.StartSpan("hash password", opentracing.ChildOf(fullSpan.Context()))
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

	err = a.streamWriter.PostRegisteredEvent(ctx, ev)

	if err != nil {
		return "", err
	}

	select {
	case <-done:
		fullSpan.LogEvent("complete")
		break

	case <-errs:
		return "", errors.Wrap(err, "failed to find registered event")

	case <-ctx.Done():
		return "", errors.New("context ended")
	}

	return UserID(id), nil
}

func (a *auth) Validate(ctx context.Context, email string, password string) (bool, error) {
	fullSpan, ctx := opentracing.StartSpanFromContextWithTracer(ctx, a.tracer, "Validate")
	defer fullSpan.Finish()

	hashSpan := a.tracer.StartSpan("hash password", opentracing.ChildOf(fullSpan.Context()))
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashSpan.Finish()

	if err != nil {
		return false, errors.Wrap(err, "failed to hash supplied password")
	}

	entry, err := a.db.GetUserByEmail(ctx, email)

	if err != nil {
		return false, errors.Wrap(err, "unexpected error while finding user")
	}

	if entry == nil {
		return false, nil
	}

	compareSpan := a.tracer.StartSpan("compare hash and password", opentracing.ChildOf(fullSpan.Context()))
	// bcrypt handles the salt in the encoded value for us
	valid := bcrypt.CompareHashAndPassword([]byte(entry.PasswordHashWithSalt), hashed) != nil
	compareSpan.Finish()

	return valid, nil
}
