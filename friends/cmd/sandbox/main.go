package main

import (
	"context"
	"log"

	"github.com/Evertras/events-demo/friends/lib/db"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d := db.New("bolt://friends-db:7687")

	err := d.Connect(ctx)

	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	defer d.Close()

	/*
		greeting, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
			result, err := transaction.Run(
				"CREATE (a:Greeting) SET a.message = $message RETURN a.message + ', from node ' + id(a)",
				map[string]interface{}{"message": "hello, world"})

				if err != nil {
					return nil, err
				}

				if result.Next() {
					return result.Record().GetByIndex(0), nil
				}

				return nil, result.Err()
		})

		if err != nil {
			log.Fatal(err)
		}

		log.Println(greeting.(string))
	*/
}
