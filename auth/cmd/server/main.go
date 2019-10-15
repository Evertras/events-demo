package main

import (
	"log"

	"github.com/Evertras/events-demo/auth/lib/authdb"
	"github.com/Evertras/events-demo/auth/lib/server"
)

const headerAuthToken = "X-Auth-Token"
const headerUserID = "X-User-ID"

func main() {
	addr := "0.0.0.0:13041"

	server := server.New(addr)

	checkDb()

	log.Println("Serving", addr)

	log.Fatal(server.ListenAndServe())
}

func checkDb() {
	db := authdb.New(authdb.ConnectionOptions{
		User:     "admin",
		Password: "admin",
		Address:  "events-demo-auth-db",
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
