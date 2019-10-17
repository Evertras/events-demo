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

	db := initDb()

	server := server.New(addr, db)

	log.Println("Serving", addr)

	log.Fatal(server.ListenAndServe())
}

func initDb() authdb.Db {
	db := authdb.New(authdb.ConnectionOptions{
		User:     "admin",
		Password: "admin",
		Address:  "auth-db",
	})

	if err := db.Connect(); err != nil {
		log.Fatalln("Error connecting to DB:", err)
	}

	log.Println("DB connected")

	if err := db.Ping(); err != nil {
		log.Fatalln("Error pinging DB:", err)
	}

	log.Println("DB pinged")

	if err := db.MigrateToLatest(); err != nil {
		log.Fatalln("Error migrating:", err)
	}

	return db
}
