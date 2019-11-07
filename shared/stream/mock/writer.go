package mock

import (
	"context"

	"github.com/Evertras/events-demo/shared/stream"
)

type StreamWrittenEvent struct {
	Key []byte
	EventID stream.EventID
	Payload stream.Serializer
}

type StreamWriter struct {
	Sent []StreamWrittenEvent
}

var _ stream.Writer = &StreamWriter{}

func NewWriter() *StreamWriter {
	return &StreamWriter{
		Sent: make([]StreamWrittenEvent, 0),
	}
}

func (s *StreamWriter) Write(
	ctx context.Context,
	key []byte,
	eventID stream.EventID,
	payload stream.Serializer,
) error {
	s.Sent = append(s.Sent, StreamWrittenEvent{
		Key: key,
		EventID: eventID,
		Payload: payload,
	})

	return nil
}

func (s *StreamWriter) Close() error {
	return nil
}
