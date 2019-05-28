package main

import "errors"

// InMemoryHivemindStore is a small in-memory implementation of a HivemindStore
type InMemoryHivemindStore struct {
	sensors map[string]Sensor
}

func (i *InMemoryHivemindStore) getSensor(id string) (Sensor, error) {
	var err error
	sensor, ok := i.sensors[id]
	if !ok {
		err = errors.New("sensor not found in store")
	}
	return sensor, err
}

func (i *InMemoryHivemindStore) getAllSensors() []Sensor {
	var sensors []Sensor
	for _, sensor := range i.sensors {
		sensors = append(sensors, sensor)
	}
	return sensors
}

func (i *InMemoryHivemindStore) storeSensorValue(id string, value Sensor) error {
	var err error
	i.sensors[id] = value
	return err
}

func (i *InMemoryHivemindStore) storeSensor(s Sensor) error {
	var err error
	i.sensors[s.ID] = s
	return err
}
