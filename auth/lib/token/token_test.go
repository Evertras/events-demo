package token

import "testing"

const testUser = "testUser"

func TestTokenCreatesAndParses(t *testing.T) {
	SignKey = []byte("abcdefkey")

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

	if parsed.UserID != testUser {
		t.Fatalf("Expected user ID %q but got %q", testUser, parsed.UserID)
	}
}
