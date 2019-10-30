package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/Evertras/events-demo/auth/lib/token"
)

func checkAuthHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span, _ := startSpan("check", r)
		defer span.Finish()

		authToken := r.Header.Get("X-Auth-Token")

		if authToken == "" {
			w.WriteHeader(401)
			log.Println("No header")
			span.SetTag("error", true)
			span.SetTag("error.object", errors.New("no header"))
			return
		}

		claim, err := token.Parse(authToken)

		if err != nil {
			w.WriteHeader(401)
			log.Println("Could not validate header:", err)
			span.SetTag("error", true)
			span.SetTag("error.object", errors.New("header validation failed"))
			return
		}

		w.Header().Set("X-User-ID", claim.Email)

		w.WriteHeader(200)
	}
}
