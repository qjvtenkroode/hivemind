package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/boltdb/bolt"
)

// integration tests
func TestIntegrationSensorAPI(t *testing.T) {
	database, err := bolt.Open("integration_test.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		t.Fatalf("setup for testing failed: %s", err)
	}
	defer database.Close()
	defer deleteDatabase(t, "integration_test.db")

	seed := []Sensor{
		Sensor{"test", 64},
	}

	err = seedBoltDB(t, database, seed)
	if err != nil {
		t.Fatalf("seed BoltDB for integration failed: %s", err)
	}

	store := BoltHivemindStore{database}
	server := NewHivemindServer(&store)

	t.Run("integration test: /api/sensor/test", func(t *testing.T) {
		want := Sensor{"test", 64}
		request := newGetRequest("api/sensor/test")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getSensorFromResponse(t, response.Body)

		assertResponseCode(t, response.Code, http.StatusOK)
		assertContentType(t, response.Header().Get("content-type"), "application/json")
		assertSensor(t, got, want)

		request = newPutRequest("api/sensor/test", strings.NewReader("12"))
		response = httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusAccepted)

		want = Sensor{"test", 12}
		request = newGetRequest("api/sensor/test")
		response = httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got = getSensorFromResponse(t, response.Body)

		assertResponseCode(t, response.Code, http.StatusOK)
		assertContentType(t, response.Header().Get("content-type"), "application/json")
		assertSensor(t, got, want)
	})

	t.Run("integration test: /api/sensor/", func(t *testing.T) {
		want := []Sensor{
			{"test", 12},
			{"third", 3},
		}

		server.ServeHTTP(httptest.NewRecorder(), newPostRequest("api/sensor/", strings.NewReader("{\"ID\": \"third\", \"Value\": 3 }")))

		request := newGetRequest("api/sensor/")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getSensorSliceFromResponse(t, response.Body)

		assertResponseCode(t, response.Code, http.StatusOK)
		assertContentType(t, response.Header().Get("content-type"), "application/json")
		assertSensorSlice(t, got, want)
	})
}
