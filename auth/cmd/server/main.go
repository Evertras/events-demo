package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Evertras/events-demo/auth/lib/authdb"
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
	w.Header().Set("X-User-ID", "user-"+authToken)

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

	checkDb()

	log.Println("Serving", addr)

	log.Fatal(server.ListenAndServe())
}

func checkDb() {
	db := authdb.New(authdb.ConnectionOptions {
		User: "admin",
		Password: "admin",
		Address: "events-demo-auth-db",
	})

	if err := db.Connect(); err != nil {
		log.Println("Error connecting to DB:", err)
		return
	}

	log.Println("DB connected")

	if err := db.Ping(); err != nil {
		log.Println("Error pinging DB:", err)
		return
	}

	log.Println("DB pinged")

	if err := db.MigrateToLatest(); err != nil {
		log.Println("Error migrating:", err)
		return
	}
}
