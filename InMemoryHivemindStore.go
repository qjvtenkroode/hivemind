package main

type InMemoryHivemindStore struct {
	sensors map[string]int
}

func (i *InMemoryHivemindStore) getSensorValue(id string) int {
	return i.sensors[id]
}
