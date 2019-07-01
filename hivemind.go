package main

// Sensor represents a sensor with an ID and current value
type Sensor struct {
	ID    string
	Name  string
	Unit  string
	Type  string
	Value int
}

// Switch represents a switch with an ID and current boolean state
type Switch struct {
	ID    string
	Name  string
	Type  string
	State bool
}

// HivemindStore is an interface for datastorage
type HivemindStore interface {
	getSensor(id string) (Sensor, error)
	getAllSensors() []Sensor
	storeSensor(s Sensor) error
	getSwitch(id string) (Switch, error)
	getAllSwitches() []Switch
	storeSwitch(s Switch) error
}
