package main

import (
	"context"
	"log"
	"net"

	"github.com/evertras/events-demo/ingest/lib/messages"
	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) Echo(ctx context.Context, in *messages.EchoRequest) (*messages.EchoReply, error) {
	log.Printf("Received: %v", in.GetMsg())

	return &messages.EchoReply{
		Reply: "Echo: " + in.GetMsg(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	messages.RegisterEchoServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
