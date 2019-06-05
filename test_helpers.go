package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/boltdb/bolt"
)

func assertBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("wrong result body; got %s, want %s", got, want)
	}
}

func assertResponseCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("wrong response code; got %d, want %d", got, want)
	}
}

func assertContentType(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("wrong content-type; got %s, want %s", got, want)
	}
}

func assertSensor(t *testing.T, got, want Sensor) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertSensorSlice(t *testing.T, got, want []Sensor) {
	t.Helper()
	sort.Slice(got, func(i, j int) bool {
		return got[i].Value > got[j].Value
	})
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func newGetRequest(url string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", url), nil)
	return req
}

func newPostRequest(url string, data io.Reader) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/%s", url), data)
	return req
}

func newPutRequest(url string, data io.Reader) *http.Request {
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/%s", url), data)
	return req
}

func randomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyz")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func getSensorFromResponse(t *testing.T, body io.Reader) (sensor Sensor) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&sensor)

	if err != nil {
		t.Fatalf("unable to parse response from server '%s' into Sensor, '%v'", body, err)
	}
	return
}

func getSensorSliceFromResponse(t *testing.T, body io.Reader) (sensors []Sensor) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&sensors)

	if err != nil {
		t.Fatalf("unable to parse response from server '%s' into []Sensor, '%v'", body, err)
	}
	return
}

func seedBoltDB(t *testing.T, database *bolt.DB, seed []Sensor) error {
	t.Helper()
	var err error

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

func deleteDatabase(t *testing.T, db string) error {
	t.Helper()
	var err error
	err = os.Remove(db)
	if err != nil {
		t.Fatalf("clean database failed: %s", err)
	}
	return err
}
