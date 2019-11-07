package stream

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/pkg/errors"

	opentracing "github.com/opentracing/opentracing-go"
	kafka "github.com/segmentio/kafka-go"
)

// Serializer is something that can Serialize to a writer, such as an Avro object
type Serializer interface {
	Serialize(buf io.Writer) error
}

type Writer interface {
	// Close gives underlying resources a chance to flush and close gracefully
	Close() error

	// Write will write a payload to the stream with the given event ID using
	// the given key to send to the correct partition.
	//
	// Note that the key is a Kafka primary key which will dictate the partition
	// and order guarantees for the event.  For example, the User ID should be used
	// on all events that have to do with a particular User.
	Write(ctx context.Context, key []byte, eventID EventID, payload Serializer) error
}

type kafkaStreamWriter struct {
	writer *kafka.Writer
}

func NewKafkaStreamWriter(topic string, brokers []string) Writer {
	cfg := kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        "user",
		Balancer:     &kafka.Hash{},
		BatchTimeout: time.Millisecond * 10,
	}

	writer := kafka.NewWriter(cfg)

	return &kafkaStreamWriter{
		writer: writer,
	}
}

func (k *kafkaStreamWriter) Write(ctx context.Context, key []byte, eventID EventID, ev Serializer) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Write")
	defer span.Finish()

	if ev == nil {
		return errors.New("nil event")
	}

	var buf bytes.Buffer

	err := ev.Serialize(&buf)

	if err != nil {
		return errors.Wrap(err, "failed to serialize event")
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: buf.Bytes(),
		Headers: []kafka.Header{
			kafka.Header{
				Key:   headerKeyEventType,
				Value: []byte(eventID),
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

