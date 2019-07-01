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
		Sensor{"test", "Test", "C", "generic", 64},
	}

	err = seedBoltDB(t, database, seed)
	if err != nil {
		t.Fatalf("seed BoltDB for integration failed: %s", err)
	}

	store := BoltHivemindStore{database}
	server := NewHivemindServer(&store)

	t.Run("integration test: /api/sensor/test", func(t *testing.T) {
		want := Sensor{"test", "Test", "C", "generic", 64}
		request := newGetRequest("api/sensor/test")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getSensorFromResponse(t, response.Body)

		assertResponseCode(t, response.Code, http.StatusOK)
		assertContentType(t, response.Header().Get("content-type"), "application/json")
		assertSensor(t, got, want)

		request = newPutRequest("api/sensor/test", strings.NewReader("{\"ID\": \"test\", \"Name\": \"Test\", \"Unit\": \"C\", \"Type\": \"generic\", \"Value\": 12}"))
		response = httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusAccepted)

		want = Sensor{"test", "Test", "C", "generic", 12}
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
			{"test", "Test", "C", "generic", 12},
			{"third", "Third", "C", "generic", 3},
		}

		server.ServeHTTP(httptest.NewRecorder(), newPostRequest("api/sensor/", strings.NewReader("{\"ID\": \"third\", \"Name\": \"Third\", \"Unit\": \"C\", \"Type\": \"generic\", \"Value\": 3 }")))

		request := newGetRequest("api/sensor/")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getSensorSliceFromResponse(t, response.Body)

		assertResponseCode(t, response.Code, http.StatusOK)
		assertContentType(t, response.Header().Get("content-type"), "application/json")
		assertSensorSlice(t, got, want)
	})
}
