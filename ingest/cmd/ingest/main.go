package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Evertras/events-demo/ingest/lib/messages"
	kafka "github.com/segmentio/kafka-go"
)

type server struct {
}

func (s *server) Echo(ctx context.Context, in *messages.EchoRequest) (*messages.EchoReply, error) {
	log.Printf("Received: %v", in.GetMsg())

	return &messages.EchoReply{
		Reply: "Echo: " + in.GetMsg(),
	}, nil
}

const deployName = "events-demo-kafka"
const kafkaService = deployName + "-cp-kafka:9092"
const zookeeperService = deployName + "-cp-zookeeper:2181"
const topic = "ingest"

var brokers = []string{kafkaService}

func main() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		ticker := time.NewTicker(time.Second)

		for {
			select {
			case <-ticker.C:
				err := writeMessages()

				if err != nil {
					log.Fatal(err)
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	/*
	go func() {
		if err := listenForMessages(ctx); err != nil {
			log.Fatal("Failed to listen: ", err)
		}
	}()
	*/

	go func() {
		<-sigchan

		log.Println("Got signal to interrupt, exiting...")

		cancel()
	}()

	<-ctx.Done()
}

func writeMessages() error {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("key-a"),
			Value: []byte("Hello"),
		},
		kafka.Message{
			Key:   []byte("key-b"),
			Value: []byte("Hello"),
		},
		kafka.Message{
			Key:   []byte("key-c"),
			Value: []byte("Hello"),
		},
		kafka.Message{
			Key:   []byte("key-d"),
			Value: []byte("Hello"),
		},
		kafka.Message{
			Key:   []byte("key-e"),
			Value: []byte("Hello"),
		},
		kafka.Message{
			Key:   []byte("key-f"),
			Value: []byte("Hello"),
		},
		kafka.Message{
			Key:   []byte("key-g"),
			Value: []byte("Hello"),
		},
		kafka.Message{
			Key:   []byte("key-h"),
			Value: []byte("Hello"),
		},
		kafka.Message{
			Key:   []byte("key-i"),
			Value: []byte("Hello"),
		},
		kafka.Message{
			Key:   []byte("key-j"),
			Value: []byte("Hello"),
		},
	)

	if err != nil {
		return err
	}

	err = w.Close()

	if err != nil {
		return err
	}

	return nil
}

func listenForMessages(ctx context.Context) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaService},
		GroupID: "abc",
		Topic:   topic,
		MaxWait: time.Second,
	})

	defer r.Close()

	log.Println("Reading...")

	go func() {
		ticker := time.NewTicker(time.Second)

		for {
			select {
			case <-ticker.C:
				_ = r.Stats()

				//log.Printf("%+v", s)

			case <-ctx.Done():
				return
			}
		}
	}()

	for {
		m, err := r.ReadMessage(ctx)

		if err != nil {
			return err
		}

		log.Printf("%d - %q %q", m.Offset, string(m.Value), m.Key)
	}
}
