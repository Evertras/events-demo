package main

import (
	"log"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/Evertras/events-demo/friends/lib/server"
	"github.com/Evertras/events-demo/shared/stream"
)

const kafkaBrokers = "kafka-cp-kafka-headless:9092"

func main() {
	addr := "0.0.0.0:13030"

	if err := initTracing(); err != nil {
		log.Fatal(err)
	}

	streamWriter := stream.NewKafkaStreamWriter("user", []string{kafkaBrokers})

	server := server.New(addr, streamWriter)

	log.Println("Serving", addr)

	log.Fatal(server.ListenAndServe())
}

func initTracing() error {
	cfg, err := jaegerconfig.FromEnv()

	if err != nil {
		return errors.Wrap(err, "failed to create tracer config")
	}

	cfg.ServiceName = "friends-api"

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
