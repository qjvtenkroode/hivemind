package main

import (
	"log"
	"net/http"
)

func main() {
	store := InMemoryHivemindStore{
		map[string]Sensor{
			"test": Sensor{"test", 5678},
		},
	}
	server := NewHivemindServer(&store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
