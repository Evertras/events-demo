package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/Evertras/events-demo/shared/stream"

	"github.com/Evertras/events-demo/friends/lib/db"
	"github.com/Evertras/events-demo/friends/lib/eventprocessor"
)

const kafkaBrokers = "kafka-cp-kafka-headless:9092"

func main() {
	ctx := context.Background()

	err := initTracing()

	if err != nil {
		log.Fatal("Failed to initialize tracing:", err)
	}

	db := initDb(ctx)

	randomID := "friends-processor-" + uuid.New().String()
	consumerGroupID, err := db.GetSharedValue(ctx, "consumerGroupID", randomID)

	if err != nil {
		log.Fatal("Failed getting consumer group ID:", err)
	}

	log.Println("Using consumer group ID", consumerGroupID)

	reader := stream.NewKafkaStreamReader("user", []string{kafkaBrokers}, consumerGroupID)
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

func initDb(ctx context.Context) db.Db {
	db := db.New("bolt://friends-db:7687")

	if err := db.Connect(ctx); err != nil {
		log.Fatalln("Error connecting to DB:", err)
	}

	log.Println("DB connected")

	return db
}

func initTracing() error {
	cfg, err := jaegerconfig.FromEnv()

	if err != nil {
		return errors.Wrap(err, "failed to create tracer config")
	}

	cfg.ServiceName = "friends-processor"

	tracer, _, err := cfg.NewTracer(
		jaegerconfig.Logger(jaegerlog.StdLogger),
		jaegerconfig.Metrics(metrics.NullFactory),
	)

	if err != nil {
		return errors.Wrap(err, "failed to create tracer")
	}

	opentracing.SetGlobalTracer(tracer)

	return nil
}
