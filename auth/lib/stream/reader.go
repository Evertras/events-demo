package stream

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"

	kafka "github.com/segmentio/kafka-go"
)

type StreamEventHandler func(data []byte) error

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

func NewKafkaStreamReader(brokers []string, groupId string) Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupId,
		Topic:   "user",
		MaxWait: time.Millisecond * 100,
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
		m, err := k.reader.ReadMessage(ctx)

		if err != nil {
			return errors.Wrap(err, "failed to read")
		}

		for _, h := range m.Headers {
			if h.Key == headerKeyEventType {
				if handler, ok := k.handlers[EventID(h.Value)]; ok {
					err = handler(m.Value)

					if err != nil {
						// TODO: Is this enough?
						log.Println("Handler error:", err)
					}
				}
			}
		}
	}
}
