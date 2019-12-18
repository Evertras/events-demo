// +build !devmode

package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// When devmode is not enabled, do actual validation
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
