package main

import (
	"testing"
	"time"

	"github.com/boltdb/bolt"
)

func TestBoltHivemindStore(t *testing.T) {
	database, err := bolt.Open("test.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		t.Fatalf("setup for testing failed: %s", err)
	}
	defer database.Close()
	defer deleteDatabase(t, "test.db")

	seed := []Sensor{
		Sensor{"13", 666},
		Sensor{"first", 1},
	}

	err = seedBoltDB(t, database, seed)
	if err != nil {
		t.Fatalf("seed BoltDB failed: %s", err)
	}

	t.Run("getSensor: json object matches", func(t *testing.T) {
		want := Sensor{"13", 666}

		store := BoltHivemindStore{database}

		got, err := store.getSensor("13")
		if err != nil {
			t.Fatalf("failure within getSensor(): %s", err)
		}

		assertSensor(t, got, want)

	})

	t.Run("getSensor: object not found", func(t *testing.T) {
		want := Sensor{}

		store := BoltHivemindStore{database}

		got, err := store.getSensor("unknown")
		if err != nil {
			t.Fatalf("failure within getSensor() for unknown: %s", err)
		}

		assertSensor(t, got, want)
	})

	t.Run("getAllSensors: get slice and match", func(t *testing.T) {
		want := []Sensor{
			Sensor{"13", 666},
			Sensor{"first", 1},
		}

		store := BoltHivemindStore{database}

		got := store.getAllSensors()

		assertSensorSlice(t, got, want)
	})

	t.Run("storeSensor: storing a new sensor", func(t *testing.T) {
		var want error
		s := Sensor{"new", 2019}

		store := BoltHivemindStore{database}

		got := store.storeSensor(s)

		if got != want {
			t.Errorf("failure within storeSensor: %s", err)
		}
	})
}
