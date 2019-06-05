package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Sensor represents a sensor with and ID and current value
type Sensor struct {
	ID    string
	Value int
}

// HivemindStore is an interface for datastorage
type HivemindStore interface {
	getSensor(id string) (Sensor, error)
	getAllSensors() []Sensor
	storeSensor(s Sensor) error
}

// HivemindServer is a HTTP interface for Hivemind
type HivemindServer struct {
	store HivemindStore
	http.Handler
}

// NewHivemindServer creates a HivemindServer with routing configured
func NewHivemindServer(s HivemindStore) *HivemindServer {
	h := new(HivemindServer)

	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(h.rootHandler))
	router.Handle("/api/", http.HandlerFunc(h.apiHandler))
	router.Handle("/api/sensor/", http.HandlerFunc(h.apiSensorHandler))

	h.Handler = router

	h.store = s

	return h
}

func (h *HivemindServer) rootHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	if url == "/" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (h *HivemindServer) apiHandler(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Path[len("/api"):]
	w.Header().Set("content-type", "application/json")
	if endpoint != "/" {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (h *HivemindServer) apiSensorHandler(w http.ResponseWriter, r *http.Request) {
	trailing := r.URL.Path[len("/api/sensor"):]
	id := strings.Split(trailing[1:], "/")[0]
	w.Header().Set("content-type", "application/json")
	switch r.Method {
	case http.MethodGet:
		h.apiSensorGet(w, trailing, id)
	case http.MethodPost:
		var body []byte
		if r.Body != nil {
			body, _ = ioutil.ReadAll(r.Body)
		}
		h.apiSensorPost(w, trailing, body)
	case http.MethodPut:
		var body []byte
		if r.Body != nil {
			body, _ = ioutil.ReadAll(r.Body)
		}
		h.apiSensorPut(w, trailing, id, body)
	}
}

func (h *HivemindServer) apiSensorGet(w http.ResponseWriter, trailing, id string) {
	if id == "" {
		err := json.NewEncoder(w).Encode(h.store.getAllSensors())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		value, err := h.store.getSensor(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		err = json.NewEncoder(w).Encode(value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (h *HivemindServer) apiSensorPost(w http.ResponseWriter, trailing string, body []byte) {
	if trailing != "/" {
		w.WriteHeader(http.StatusNotImplemented)
	} else {
		var s Sensor
		err := json.Unmarshal(body, &s)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = h.store.storeSensor(s)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

func (h *HivemindServer) apiSensorPut(w http.ResponseWriter, trailing, id string, body []byte) {
	value, _ := strconv.Atoi(string(body))
	err := h.store.storeSensor(Sensor{id, value})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
