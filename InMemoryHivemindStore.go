package main

// InMemoryHivemindStore is a small in-memory implementation of a HivemindStore
type InMemoryHivemindStore struct {
	sensors map[string]int
}

func (i *InMemoryHivemindStore) getSensorValue(id string) int {
	return i.sensors[id]
}

func (i *InMemoryHivemindStore) storeSensorValue(id string, value int) error {
	var err error
	i.sensors[id] = value
	return err
}
