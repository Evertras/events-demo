package auth

import (
	"testing"

	"github.com/Evertras/events-demo/auth/lib/authdb/mockauthdb"
)

func TestRegisterUser(t *testing.T) {
	db := mockauthdb.New()
	a := New(db)

	email := "test@testing.com"
	password := "sekrit"
	details := RegistrationMeta{}

	err := a.Register(email, password, details)

	if err != nil {
		t.Fatal("Failed to register", err)
	}

	entry, exists := db.EntriesByEmail[email]

	if !exists {
		t.Fatal("Entry not found in DB")
		return
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.f)
	}
}
