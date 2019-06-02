package main

import (
	"encoding/json"
	"os"
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
	defer deleteDatabase(t)

	err = seedBoltDB(t, database)
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

	t.Run("storeSensorValue: storing a new value", func(t *testing.T) {
		var want error
		s := Sensor{"13", 1988}

		store := BoltHivemindStore{database}

		got := store.storeSensorValue(s.ID, s)

		if got != want {
			t.Errorf("failure within storeSensorValue: %s", err)
		}
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

// helpers
func seedBoltDB(t *testing.T, database *bolt.DB) error {
	t.Helper()
	var err error

	seed := []Sensor{
		Sensor{"13", 666},
		Sensor{"first", 1},
	}

	for _, s := range seed {
		encoded, err := json.Marshal(s)
		if err != nil {
			return err
		}
		err = database.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte("sensor"))
			if err != nil {
				return err
			}

			err = bucket.Put([]byte(s.ID), encoded)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			t.Fatalf("seed failed: %s", err)
		}
	}

	return err
}

func deleteDatabase(t *testing.T) error {
	t.Helper()
	var err error
	err = os.Remove("test.db")
	if err != nil {
		t.Fatalf("clean database failed: %s", err)
	}
	return err
}
