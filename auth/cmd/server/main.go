package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"

	"github.com/google/uuid"

	"github.com/Evertras/events-demo/shared/stream"
	"github.com/Evertras/events-demo/auth/lib/auth"
	"github.com/Evertras/events-demo/auth/lib/authdb"
	"github.com/Evertras/events-demo/auth/lib/events"
	"github.com/Evertras/events-demo/auth/lib/server"
	"github.com/Evertras/events-demo/auth/lib/token"
	"github.com/Evertras/events-demo/auth/lib/tracing"
)

const headerAuthToken = "X-Auth-Token"
const headerUserID = "X-User-ID"
const addr = "0.0.0.0:13041"

const kafkaBrokers = "kafka-cp-kafka-headless:9092"

func main() {
	ctx := context.Background()

	err := tracing.Init("auth-api")

	db := initDb(ctx)

	err = initSignKey(ctx, db)

	if err != nil {
		log.Fatal("Failed to initialize token sign key:", err)
	}

	streamWriter := stream.NewKafkaStreamWriter("user", []string{kafkaBrokers})
	writer := events.NewWriter(streamWriter)
	a := auth.New(db, writer)
	server := server.New(addr, a)

	if err != nil {
		log.Fatal("Failed to create server:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	log.Println("Serving", addr)

	log.Fatal(server.ListenAndServe())
}

func initDb(ctx context.Context) authdb.Db {
	db := authdb.New(authdb.ConnectionOptions{
		Address: "auth-db:6379",
	})

	if err := db.Connect(ctx); err != nil {
		log.Fatalln("Error connecting to DB:", err)
	}

	log.Println("DB connected")

	if err := db.Ping(ctx); err != nil {
		log.Fatalln("Error pinging DB:", err)
	}

	log.Println("DB pinged")

	return db
}

func initSignKey(ctx context.Context, db authdb.Db) error {
	buf := make([]byte, 1024)

	rand.Reader.Read(buf)

	randomSignKey := base64.StdEncoding.EncodeToString(buf)

	tokenSignKey, err := db.GetSharedValue(ctx, "auth.token.signKey", randomSignKey)

	if err != nil {
		return err
	}

	token.SignKey, err = base64.StdEncoding.DecodeString(tokenSignKey)

	return err
}
