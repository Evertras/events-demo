package stream

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/pkg/errors"

	kafka "github.com/segmentio/kafka-go"

	"github.com/Evertras/events-demo/auth/lib/stream/authevents"
)


// EventStream represents an event stream such as Kafka or Redis pub/sub
type Writer interface {
	// PostRegisteredEvent posts a UserRegistered event to the stream
	PostRegisteredEvent(ctx context.Context, ev *authevents.UserRegistered) error

	// Close gives underlying resources a chance to flush and close gracefully
	Close() error
}

type kafkaStreamWriter struct {
	writer *kafka.Writer
}

func NewKafkaStreamWriter(brokers []string) Writer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        "user",
		Balancer:     &kafka.Hash{},
		BatchTimeout: time.Millisecond * 10,
	})

	return &kafkaStreamWriter{
		writer: writer,
	}
}

func (k *kafkaStreamWriter) PostRegisteredEvent(ctx context.Context, ev *authevents.UserRegistered) error {
	return k.write(ctx, []byte(ev.ID), ev)
}

type serializer interface {
	Serialize(buf io.Writer) error
}

func (k *kafkaStreamWriter) write(ctx context.Context, key []byte, ev serializer) error {
	if ev == nil {
		return errors.New("nil event")
	}

	var buf bytes.Buffer

	ev.Serialize(&buf)

	err := k.writer.WriteMessages(
		ctx,
		kafka.Message{
			Key:   key,
			Value: buf.Bytes(),
			Headers: []kafka.Header {
				kafka.Header{
					Key: headerKeyEventType,
					Value: []byte("UserRegistered"),
				},
			},
		})

	if err != nil {
		return errors.Wrap(err, "failed to write messages")
	}

	return nil
}

func (k *kafkaStreamWriter) Close() error {
	return k.writer.Close()
}
