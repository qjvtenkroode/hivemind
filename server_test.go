package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
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

// helpers
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
