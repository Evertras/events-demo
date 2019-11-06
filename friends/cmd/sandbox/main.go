package main

import (
	"log"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func main() {

	driver, err := neo4j.NewDriver("bolt://friends-db:7687", neo4j.NoAuth())
	if err != nil {
		log.Fatal("Failed to create driver:", err)
	}
	defer driver.Close()

	log.Println("Driver created")

	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		log.Fatal("Failed to create session:", err)
	}
	defer session.Close()

	log.Println("Session created")

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
}
