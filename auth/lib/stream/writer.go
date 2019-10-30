package stream

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/pkg/errors"

	opentracing "github.com/opentracing/opentracing-go"
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
	cfg := kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        "user",
		Balancer:     &kafka.Hash{},
		BatchTimeout: time.Millisecond * 10,
		// Logger:   log.New(os.Stdout, "kafka-writer ", log.LstdFlags),
	}

	writer := kafka.NewWriter(cfg)

	return &kafkaStreamWriter{
		writer: writer,
	}
}

func (k *kafkaStreamWriter) PostRegisteredEvent(ctx context.Context, ev *authevents.UserRegistered) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostRegisteredEvent")
	defer span.Finish()

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

	err := ev.Serialize(&buf)

	if err != nil {
		return errors.Wrap(err, "failed to serialize event")
	}

	msg := kafka.Message{
		Key:   key,
		Value: buf.Bytes(),
		Headers: []kafka.Header{
			kafka.Header{
				Key:   headerKeyEventType,
				Value: []byte("UserRegistered"),
			},
		},
	}

	if span := opentracing.SpanFromContext(ctx); span != nil {
		spanCtx := span.Context()

		var ctxBuf bytes.Buffer

		err := opentracing.GlobalTracer().Inject(spanCtx, opentracing.Binary, &ctxBuf)

		if err != nil {
			return errors.Wrap(err, "failed to inject span context")
		}

		msg.Headers = append(msg.Headers,
			kafka.Header{
				Key:   headerKeySpanContext,
				Value: ctxBuf.Bytes(),
			},
		)
	}

	err = k.writer.WriteMessages(
		ctx,
		msg,
	)

	if err != nil {
		return errors.Wrap(err, "failed to write messages")
	}

	return nil
}

func (k *kafkaStreamWriter) Close() error {
	return k.writer.Close()
}
