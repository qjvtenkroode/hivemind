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

func (b *BoltHivemindStore) getSwitch(id string) (Switch, error) {
	var err error
	var sw Switch
	err = b.database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("switch"))
		if bucket == nil {
			return nil
		}
		v := bucket.Get([]byte(id))
		if v == nil {
			return nil
		}
		err = json.Unmarshal(v, &sw)
		if err != nil {
			return err
		}
		return nil
	})

	return sw, err
}

func (b *BoltHivemindStore) getAllSwitches() []Switch {
	var switches []Switch

	_ = b.database.View(func(tx *bolt.Tx) error {
		var sw Switch
		bucket := tx.Bucket([]byte("switch"))
		if bucket == nil {
			return nil
		}
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			err := json.Unmarshal(v, &sw)
			if err != nil {
				return err
			}
			switches = append(switches, sw)
		}
		return nil
	})

	return switches
}

func (b *BoltHivemindStore) storeSwitch(sw Switch) error {
	var err error

	err = b.database.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("switch"))
		if err != nil {
			return err
		}
		encoded, err := json.Marshal(sw)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(sw.ID), encoded)
	})

	return err
}
