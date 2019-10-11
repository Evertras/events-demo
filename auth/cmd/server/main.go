package main

import (
	"log"
	"net/http"
	"time"
)

const headerAuthToken = "X-Auth-Token"
const headerUserID = "X-User-ID"

func checkAuthHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("X-Auth-Token")

	if authToken == "" {
		w.WriteHeader(401)
		log.Println("unauthorized")
		return
	}

	log.Println("OK!", authToken)

	// Right now everyone can get by as long as they have any token, super secure!
	w.Header().Set("X-User-ID", "user-" + authToken)

	w.WriteHeader(200)
}

func main() {
	addr := "0.0.0.0:13041"

	router := http.NewServeMux()

	router.HandleFunc("/check", checkAuthHandler)

	server := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		Handler:      router,
	}

	log.Println("Serving", addr)

	log.Fatal(server.ListenAndServe())
}
