package token

import "testing"

const testUser = "testUser"

func TestTokenCreatesAndParses(t *testing.T) {
	token, err := New(testUser)

	if err != nil {
		t.Fatal(err)
	}

	parsed, err := Parse(token)

	if err != nil {
		t.Fatal(err)
	}

	if err = parsed.Valid(); err != nil {
		t.Fatal(err)
	}

	if parsed.Email != testUser {
		t.Fatalf("Expected username %q but got %q", testUser, parsed.Email)
	}
}
