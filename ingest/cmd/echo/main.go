package main

import (
	"context"
	"log"
	"time"

	"github.com/evertras/events-demo/ingest/lib/messages"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not dial: %v", err)
	}

	defer conn.Close()

	client := messages.NewEchoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reply, err := client.Echo(ctx, &messages.EchoRequest{ Msg: "Hello!" })

	log.Print(reply)
}
