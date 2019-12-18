// +build devmode

package auth

import (
	"context"
)

// Inject a backdoor if we build with devmode enabled
func (a *auth) ValidateByID(ctx context.Context, id string, password string) (bool, error) {
	return true, nil
}
