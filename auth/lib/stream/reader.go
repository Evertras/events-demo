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

	"github.com/Evertras/events-demo/auth/lib/tracing"
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
	tracer    opentracing.Tracer
}

func NewKafkaStreamReader(brokers []string, groupId string) (Reader, error) {
	tracer, err := tracing.Init("kafka-reader")

	if err != nil {
		return nil, errors.Wrap(err, "failed to init tracer:")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupId,
		Topic:       "user",
		MaxWait:     time.Millisecond * 100,
		MaxAttempts: 5,
	})

	return &kafkaStreamReader{
		reader:    reader,
		handlers:  make(map[EventID]StreamEventHandler),
		listening: false,
		tracer:    tracer,
	}, nil
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
		m, err := k.reader.ReadMessage(ctx)

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
				spanCtx, err = k.tracer.Extract(opentracing.Binary, buf)

				if err != nil {
					log.Println("Error getting span context:", err)
				}
			}
		}

		if evType != "" {
			if handler, ok := k.handlers[evType]; ok {
				span := k.tracer.StartSpan("process "+string(evType), ext.RPCServerOption(spanCtx))

				err = handler(opentracing.ContextWithSpan(ctx, span), m.Value)

				if err != nil {
					// TODO: Is this enough?
					log.Println("Handler error:", err)
				}

				span.Finish()
			}
		}
	}
}
