package main

import (
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
)

func main() {
	database, err := bolt.Open("hivemind.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalf("setup of hivemind.db failed: %s", err)
	}
	defer database.Close()

	store := BoltHivemindStore{database}
	server := NewHivemindServer(&store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
