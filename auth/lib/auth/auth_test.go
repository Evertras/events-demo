package auth

import (
	"testing"

	"github.com/Evertras/events-demo/auth/lib/authdb/mockauthdb"
	"github.com/Evertras/events-demo/auth/lib/eventstream/mockeventstream"
	"github.com/Evertras/events-demo/auth/lib/eventstream/events"
)

func TestRegisterUser(t *testing.T) {
	db := mockauthdb.New()
	stream := mockeventstream.New()

	a := New(db, stream)

	email := "test@testing.com"
	password := "sekrit"
	username := "some-user"
	details := RegistrationMeta{
		Username: username,
	}

	err := a.Register(email, password, details)

	if err != nil {
		t.Fatal("Failed to register", err)
	}

	entry, exists := db.EntriesByEmail[email]

	if !exists {
		t.Fatal("Entry not found in DB")
	}

	testCases := []struct{
		name string
		f func(t *testing.T)
	}{
		{
			name: "PasswordHashed",
			f: func(t *testing.T) {
				if entry.PasswordHashWithSalt == password {
					t.Error("Password is in plaintext")
				}
			},
		},
		{
			name: "GeneratesID",
			f: func(t *testing.T) {
				if len(entry.ID) == 0 {
					t.Error("ID not generated")
				}
			},
		},
		{
			name: "RegisteredEventSent",
			f: func(t *testing.T) {
				if len(stream.Sent) != 1 {
					t.Fatalf("Expected 1 event sent, but found %d", len(stream.Sent))
				}

				ev := stream.Sent[0].(*events.UserRegistered)

				if ev == nil {
					t.Fatal("Expected events.UserRegistered but did not cast correctly")
				}

				if ev.ID != entry.ID {
					t.Errorf("Did not send correct ID: expected %q but got %q", entry.ID, ev.ID)
				}

				if ev.Email != email {
					t.Errorf("Did not send correct user email: expected %q but got %q", email, ev.Email)
				}

				if ev.Username != username {
					t.Errorf("Did not send correct username: expected %q but got %q", username, ev.Username)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.f)
	}
}
