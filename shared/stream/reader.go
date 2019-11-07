package stream

import (
	"bytes"
	"context"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	kafka "github.com/segmentio/kafka-go"
)

type StreamEventHandler func(ctx context.Context, data []byte) error

type Reader interface {
	// RegisterHandler adds an event handler that will perform some action
	// when an event with the given ID is received.  Must be called before
	// calling Listen.
	RegisterHandler(event EventID, handler StreamEventHandler) error

	// Listen will listen for incoming events and route them internally.
	// This is a blocking call.
	Listen(ctx context.Context) error
}

type kafkaStreamReader struct {
	reader    *kafka.Reader
	handlers  map[EventID]StreamEventHandler
	listening bool
	lock      sync.Mutex
}

func NewKafkaStreamReader(topic string, brokers []string, groupId string) Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:               brokers,
		GroupID:               groupId,
		Topic:                 topic,
		WatchPartitionChanges: true,
		MaxWait:               time.Millisecond * 10,
		// Logger:                log.New(os.Stdout, "kafka-reader ", log.LstdFlags),
	})

	return &kafkaStreamReader{
		reader:    reader,
		handlers:  make(map[EventID]StreamEventHandler),
		listening: false,
	}
}

func (k *kafkaStreamReader) RegisterHandler(event EventID, handler StreamEventHandler) error {
	k.lock.Lock()
	defer k.lock.Unlock()

	if k.listening {
		return errors.New("cannot add a handler when already listening")
	}

	k.handlers[event] = handler

	return nil
}

func (k *kafkaStreamReader) Listen(ctx context.Context) error {
	k.lock.Lock()
	k.listening = true
	k.lock.Unlock()

	defer func() {
		k.lock.Lock()
		k.listening = false
		k.lock.Unlock()
	}()

	for {
		m, err := k.reader.FetchMessage(ctx)

		// This is fatal
		if err != nil {
			return errors.Wrap(err, "failed to read")
		}

		var evType EventID
		var spanCtx opentracing.SpanContext

		for _, h := range m.Headers {
			switch h.Key {
			case headerKeyEventType:
				evType = EventID(h.Value)

			case headerKeySpanContext:
				buf := bytes.NewBuffer(h.Value)
				spanCtx, err = opentracing.GlobalTracer().Extract(opentracing.Binary, buf)

				if err != nil {
					log.Println("Error getting span context:", err)
				}
			}
		}

		if evType != "" {
			if handler, ok := k.handlers[evType]; ok {
				span := opentracing.StartSpan("process "+string(evType), ext.RPCServerOption(spanCtx))

				err = handler(opentracing.ContextWithSpan(ctx, span), m.Value)

				// This is not fatal... for now
				if err != nil {
					log.Println("Handler error:", err)

					span.SetTag("error", true)
					span.SetTag("error.object", err)
				}

				span.Finish()
			}
		}

		err = k.reader.CommitMessages(ctx, m)

		// This is fatal
		if err != nil {
			return errors.Wrap(err, "failed to commit messages")
		}
	}
}
