package main

import (
	"context"
	"log"

	"github.com/Evertras/events-demo/auth/lib/auth"
	"github.com/Evertras/events-demo/auth/lib/authdb"
	"github.com/Evertras/events-demo/auth/lib/eventprocessor"
	"github.com/Evertras/events-demo/auth/lib/server"
	"github.com/Evertras/events-demo/auth/lib/stream"
)

const headerAuthToken = "X-Auth-Token"
const headerUserID = "X-User-ID"

const kafkaBrokers = "kafka-cp-kafka-headless:9092"

func main() {
	addr := "0.0.0.0:13041"

	db := initDb()

	consumerGroupID, err := db.GetSharedID("auth.consumerGroupID")

	if err != nil {
		log.Fatal("Failed getting consumer group ID", err)
	}

	log.Println("Using consumer group ID", consumerGroupID)

	writer := initStreamWriter()
	reader := initStreamReader(consumerGroupID)

	a := auth.New(db, writer)

	server := server.New(addr, a)

	processor := eventprocessor.New(db, reader)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	go func() {
		err := reader.Listen(ctx)

		if err != nil {
			log.Fatalln("Error listening:", err)
		}
	}()

	go func() {
		err := processor.Run(ctx)
		log.Println("Processor finished")
		if err != nil {
			log.Fatalln("Error processing:", err)
		}
	}()

	log.Println("Serving", addr)

	log.Fatal(server.ListenAndServe())
}

func initDb() authdb.Db {
	db := authdb.New(authdb.ConnectionOptions{
		Address: "auth-db:6379",
	})

	if err := db.Connect(); err != nil {
		log.Fatalln("Error connecting to DB:", err)
	}

	log.Println("DB connected")

	if err := db.Ping(); err != nil {
		log.Fatalln("Error pinging DB:", err)
	}

	log.Println("DB pinged")

	return db
}

func initStreamWriter() stream.Writer {
	return stream.NewKafkaStreamWriter([]string{kafkaBrokers})
}

func initStreamReader(groupId string) stream.Reader {
	return stream.NewKafkaStreamReader([]string{kafkaBrokers}, groupId)
}
