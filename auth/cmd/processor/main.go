package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/google/uuid"

	"github.com/Evertras/events-demo/auth/lib/authdb"
	"github.com/Evertras/events-demo/auth/lib/eventprocessor"
	"github.com/Evertras/events-demo/auth/lib/stream"
	"github.com/Evertras/events-demo/auth/lib/tracing"
)

const kafkaBrokers = "kafka-cp-kafka-headless:9092"

func main() {
	ctx := context.Background()

	err := tracing.Init("auth-processor")

	db := initDb(ctx)

	randomID := "auth-processor-" + uuid.New().String()
	consumerGroupID, err := db.GetSharedValue(ctx, "auth.consumerGroupID", randomID)

	if err != nil {
		log.Fatal("Failed getting consumer group ID:", err)
	}

	log.Println("Using consumer group ID", consumerGroupID)

	reader := stream.NewKafkaStreamReader([]string{kafkaBrokers}, consumerGroupID)
	processor := eventprocessor.New(db)

	processor.RegisterHandlers(reader)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	go func() {
		err := reader.Listen(ctx)

		if err != nil {
			log.Fatalln("Error listening:", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt)

	log.Println("Processing...")

	<-signalChan

	log.Println("Closed gracefully")
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
