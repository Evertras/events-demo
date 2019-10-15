package server

import (
	"log"
	"net/http"
)

func checkAuthHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("X-Auth-Token")

	if authToken == "" {
		w.WriteHeader(401)
		log.Println("unauthorized")
		return
	}

	log.Println("OK!", authToken)

	// Right now everyone can get by as long as they have any token, super secure!
	w.Header().Set("X-User-ID", "user-"+authToken)

	w.WriteHeader(200)
}
