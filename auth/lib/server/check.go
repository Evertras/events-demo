package server

import (
	"log"
	"net/http"

	"github.com/Evertras/events-demo/auth/lib/token"
)

func checkAuthHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("X-Auth-Token")

	if authToken == "" {
		w.WriteHeader(401)
		log.Println("No header")
		return
	}

	claim, err := token.Parse(authToken)

	if err != nil {
		w.WriteHeader(401)
		log.Println("Could not validate header:", err)
		return
	}

	w.Header().Set("X-User-ID", claim.Username)

	w.WriteHeader(200)
}
