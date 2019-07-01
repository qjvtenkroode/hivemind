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
			"test":   Sensor{"test", "Test", "C", "generic", 64},
			"second": Sensor{"second", "Second", "C", "generic", 2},
		},
		nil,
	}
	server := NewHivemindServer(&store)

	t.Run("return json value: 64, status 200 on GET /api/sensor/test", func(t *testing.T) {
		want := Sensor{"test", "Test", "C", "generic", 64}
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
			{"test", "Test", "C", "generic", 64},
			{"second", "Second", "C", "generic", 2},
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
		request := newPostRequest("api/sensor/", strings.NewReader("{\"ID\": \"status_202\", \"Name\": \"Status 202\", \"Unit\": \"C\", \"Type\": \"generic\", \"Value\": 202}"))
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
		request := newPutRequest("api/sensor/test", strings.NewReader("{\"ID\": \"test\", \"Name\": \"Test\", \"Unit\": \"C\", \"Type\": \"generic\", \"Value\": 1234}"))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusAccepted)
	})
}

func TestSwitchAPI(t *testing.T) {
	store := StubHivemindStore{
		nil,
		map[string]Switch{
			"test":   Switch{"test", "test", "generic", true},
			"second": Switch{"second", "second", "generic", false},
		},
	}
	server := NewHivemindServer(&store)

	t.Run("return json value: true, status 200 on GET /api/switch/test", func(t *testing.T) {
		want := Switch{"test", "test", "generic", true}
		request := newGetRequest("api/switch/test")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getSwitchFromResponse(t, response.Body)

		assertResponseCode(t, response.Code, http.StatusOK)
		assertContentType(t, response.Header().Get("content-type"), "application/json")
		assertSwitch(t, got, want)
	})

	t.Run("return status 404 on GET /api/switch/{random}", func(t *testing.T) {
		request := newGetRequest(fmt.Sprintf("api/switch/%s", randomString(8)))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusNotFound)
	})

	t.Run("return api switch table as json, status 200 on GET /api/sensor/", func(t *testing.T) {
		want := []Switch{
			{"test", "test", "generic", true},
			{"second", "second", "generic", false},
		}

		request := newGetRequest("api/switch/")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getSwitchSliceFromResponse(t, response.Body)

		assertResponseCode(t, response.Code, http.StatusOK)
		assertContentType(t, response.Header().Get("content-type"), "application/json")
		assertSwitchSlice(t, got, want)
	})

	t.Run("return status 202 on POST /api/switch/", func(t *testing.T) {
		request := newPostRequest("api/switch/", strings.NewReader("{\"ID\": \"status_202\", \"Name\": \"Status 200\", \"State\": false}"))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusAccepted)
	})

	t.Run("return status 501 on POST /api/switch/test", func(t *testing.T) {
		request := newPostRequest("api/switch/test", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusNotImplemented)
	})

	t.Run("return status 202 on PUT /api/switch/test", func(t *testing.T) {
		request := newPutRequest("api/switch/test", strings.NewReader("{\"ID\": \"test\", \"Name\": \"test\", \"State\": false}"))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusAccepted)
	})
}

// stubs
type StubHivemindStore struct {
	sensors  map[string]Sensor
	switches map[string]Switch
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

func (s *StubHivemindStore) storeSensor(sensor Sensor) error {
	var err error
	s.sensors[sensor.ID] = sensor
	return err
}

func (s *StubHivemindStore) getSwitch(id string) (Switch, error) {
	var err error
	sw, ok := s.switches[id]
	if !ok {
		err = errors.New("Switch not found in store")
	}
	return sw, err
}

func (s *StubHivemindStore) getAllSwitches() []Switch {
	var switches []Switch
	for _, sw := range s.switches {
		switches = append(switches, sw)
	}
	return switches
}

func (s *StubHivemindStore) storeSwitch(sw Switch) error {
	var err error
	s.switches[sw.ID] = sw
	return err
}
