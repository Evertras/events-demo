package eventstream

import (
	"github.com/pkg/errors"

	kafka "github.com/segmentio/kafka-go"

	"github.com/Evertras/events-demo/auth/lib/eventstream/events"
)

type EventStream interface {
	PostRegisteredEvent(ev *events.UserRegistered) error
}

type kafkaStream struct {
	writer *kafka.Writer
}

func NewKafkaStream(brokers []string) EventStream {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    "user",
		Balancer: &kafka.LeastBytes{},
	})

	return &kafkaStream{
		writer: writer,
	}
}

func (k *kafkaStream) PostRegisteredEvent(ev *events.UserRegistered) error {
	return errors.New("not implemented")
}
