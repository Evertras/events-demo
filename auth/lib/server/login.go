package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Evertras/events-demo/auth/lib/auth"
	"github.com/Evertras/events-demo/auth/lib/token"
)

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func loginHandler(auth auth.Auth) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span, ctx := startSpan("login", r)
		defer span.Finish()

		if r.Method != "POST" {
			w.WriteHeader(400)
			log.Println("Method must be POST")
			span.SetTag("error", true)
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to parse:", err)
			span.SetTag("error", true)
			return
		}

		var login LoginBody
		err = json.Unmarshal(body, &login)

		if err != nil {
			w.WriteHeader(400)
			log.Println("Could not parse body:", err)
			span.SetTag("error", true)
			return
		}

		id, err := auth.GetIDFromEmail(ctx, login.Email)

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to get ID:", err)
			span.SetTag("error", true)
			return
		}

		valid, err := auth.ValidateByID(ctx, id, login.Password)

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to validate user:", err)
			span.SetTag("error", true)
			return
		}

		if !valid {
			w.WriteHeader(400)
			log.Println("Bad credentials for", login.Email)
			span.SetTag("error", true)
			return
		}

		t, err := token.New(id)

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to generate token:", err)
			span.SetTag("error", true)
			return
		}

		resBody, err := json.Marshal(TokenResponse{Token: t})

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to marshal token response:", err)
			span.SetTag("error", true)
			return
		}

		_, err = w.Write(resBody)

		if err != nil {
			w.WriteHeader(500)
			log.Println("Failed to write body:", err)
			span.SetTag("error", true)
			return
		}

		log.Println("Login successful for " + login.Email)
	}
}
