package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	server := NewHivemindServer(nil)

	t.Run("return empty body and status 200 on /", func(t *testing.T) {
		request := newGetRequest("")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertBody(t, response.Body.String(), "")
		assertResponseCode(t, response.Code, http.StatusOK)
	})

	t.Run("return empty body and status 404 on /unknown", func(t *testing.T) {
		request := newGetRequest("unknown")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertBody(t, response.Body.String(), "")
		assertResponseCode(t, response.Code, http.StatusNotFound)
	})

	t.Run("return empty body and status 200 on /api", func(t *testing.T) {
		request := newGetRequest("api/")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertBody(t, response.Body.String(), "")
		assertResponseCode(t, response.Code, http.StatusOK)
		assertContentType(t, response.Header().Get("content-type"), "application/json")
	})

	t.Run("return status 501 on /api/{random}", func(t *testing.T) {
		request := newGetRequest(fmt.Sprintf("api/%s", randomString(8)))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusNotImplemented)
	})
}

func TestSensorAPI(t *testing.T) {
	store := StubHivemindStore{
		map[string]Sensor{
			"test":   Sensor{"test", 64},
			"second": Sensor{"second", 2},
		},
	}
	server := NewHivemindServer(&store)

	t.Run("return json value: 64, status 200 on GET /api/sensor/test", func(t *testing.T) {
		want := Sensor{"test", 64}
		request := newGetRequest("api/sensor/test")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getSensorFromResponse(t, response.Body)

		assertResponseCode(t, response.Code, http.StatusOK)
		assertContentType(t, response.Header().Get("content-type"), "application/json")
		assertSensor(t, got, want)
	})

	t.Run("return status 404 on GET /api/sensor/{random}", func(t *testing.T) {
		request := newGetRequest(fmt.Sprintf("api/sensor/%s", randomString(8)))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusNotFound)
	})

	t.Run("return api sensor table as json, status 200 on GET /api/sensor/", func(t *testing.T) {
		want := []Sensor{
			{"test", 64},
			{"second", 2},
		}

		request := newGetRequest("api/sensor/")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getSensorSliceFromResponse(t, response.Body)

		assertResponseCode(t, response.Code, http.StatusOK)
		assertContentType(t, response.Header().Get("content-type"), "application/json")
		assertSensorSlice(t, got, want)
	})

	t.Run("return status 202 on POST /api/sensor/", func(t *testing.T) {
		request := newPostRequest("api/sensor/", strings.NewReader("{\"ID\": \"status_202\", \"Value\": 202}"))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusAccepted)
	})

	t.Run("return status 501 on POST /api/sensor/test", func(t *testing.T) {
		request := newPostRequest("api/sensor/test", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusNotImplemented)
	})

	t.Run("return status 202 on PUT /api/sensor/test", func(t *testing.T) {
		request := newPutRequest("api/sensor/test", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusAccepted)
	})
}

// stubs
type StubHivemindStore struct {
	sensors map[string]Sensor
}

func (s *StubHivemindStore) getSensor(id string) (Sensor, error) {
	var err error
	sensor, ok := s.sensors[id]
	if !ok {
		err = errors.New("Sensor not found in store")
	}
	return sensor, err
}

func (s *StubHivemindStore) getAllSensors() []Sensor {
	var sensors []Sensor
	for _, sensor := range s.sensors {
		sensors = append(sensors, sensor)
	}
	return sensors
}

func (s *StubHivemindStore) storeSensorValue(id string, value Sensor) error {
	var err error
	s.sensors[id] = value
	return err
}

func (s *StubHivemindStore) storeSensor(sensor Sensor) error {
	var err error
	s.sensors[sensor.ID] = sensor
	return err
}
