package main

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

// BoltHivemindStore is a HivemindStore implementation based on BoltDB
type BoltHivemindStore struct {
	database *bolt.DB
}

func (b *BoltHivemindStore) getSensor(id string) (Sensor, error) {
	var err error
	var sensor Sensor
	err = b.database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("sensor"))
		if bucket == nil {
			return nil
		}
		v := bucket.Get([]byte(id))
		if v == nil {
			return nil
		}
		err = json.Unmarshal(v, &sensor)
		if err != nil {
			return err
		}
		return nil
	})

	return sensor, err
}

func (b *BoltHivemindStore) getAllSensors() []Sensor {
	var sensors []Sensor

	_ = b.database.View(func(tx *bolt.Tx) error {
		var s Sensor
		bucket := tx.Bucket([]byte("sensor"))
		if bucket == nil {
			return nil
		}
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			err := json.Unmarshal(v, &s)
			if err != nil {
				return err
			}
			sensors = append(sensors, s)
		}
		return nil
	})

	return sensors
}

func (b *BoltHivemindStore) storeSensorValue(id string, value Sensor) error {
	var err error

	err = b.database.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("sensor"))
		if err != nil {
			return err
		}
		encoded, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(id), encoded)
	})

	return err
}

func (b *BoltHivemindStore) storeSensor(sensor Sensor) error {
	var err error

	err = b.database.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("sensor"))
		if err != nil {
			return err
		}
		encoded, err := json.Marshal(sensor)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(sensor.ID), encoded)
	})

	return err
}
