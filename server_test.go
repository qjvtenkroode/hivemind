package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	server := NewHivemindServer()

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
	server := NewHivemindServer()

	t.Run("return json value: 64, status 200 on GET /api/sensor/test", func(t *testing.T) {
		request := newGetRequest("api/sensor/test")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertBody(t, response.Body.String(), "64")
		assertResponseCode(t, response.Code, http.StatusOK)
		assertContentType(t, response.Header().Get("content-type"), "application/json")
	})
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

func newGetRequest(url string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", url), nil)
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
