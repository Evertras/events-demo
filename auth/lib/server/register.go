package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Evertras/events-demo/auth/lib/auth"
	"github.com/Evertras/events-demo/auth/lib/token"
)

type RegisterBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func registerHandler(a auth.Auth) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span, ctx := startSpan("register", r)
		defer span.Finish()

		if r.Method != "POST" {
			w.WriteHeader(400)
			log.Println("Method must be POST")
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to parse:", err)
			return
		}

		var register RegisterBody
		err = json.Unmarshal(body, &register)

		if err != nil {
			w.WriteHeader(400)
			log.Println("Could not parse body:", err)
			return
		}

		if register.Email == "" {
			w.WriteHeader(400)
			log.Println("Missing email address")
			return
		}

		if register.Password == "" {
			w.WriteHeader(400)
			log.Println("Missing password")
			return
		}

		_, err = a.Register(ctx, register.Email, register.Password)

		if err != nil {
			if err == auth.ErrUserAlreadyExists {
				w.WriteHeader(400)
				log.Println("User already exists:", register.Email)
				return
			}
			w.WriteHeader(500)
			log.Println("Failed to register user:", err)
			return
		}

		valid, err := a.Validate(ctx, register.Email, register.Password)

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to validate:", err)
			return
		}

		if !valid {
			w.WriteHeader(400)
			log.Println("Failed to validate credentials for", register.Email)
			return
		}

		t, err := token.New(register.Email)

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to generate token:", err)
			return
		}

		resBody, err := json.Marshal(TokenResponse{Token: t})

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to marshal token response:", err)
			return
		}

		_, err = w.Write(resBody)

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to write body:", err)
			return
		}
	}
}
