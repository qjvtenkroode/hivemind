package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// HivemindStore is an interface for datastorage
type HivemindStore interface {
	getSensorValue(id string) int
	storeSensorValue(id string, value int) error
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
		h.apiSensorPost(w, trailing, id)
	case http.MethodPut:
		var body []byte
		if r.Body != nil {
			body, _ = ioutil.ReadAll(r.Body)
		}
		h.apiSensorPut(w, trailing, id, body)
	}
}

func (h *HivemindServer) apiSensorGet(w http.ResponseWriter, trailing, id string) {
	value := h.store.getSensorValue(id)
	if value == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, value)
}

func (h *HivemindServer) apiSensorPost(w http.ResponseWriter, trailing, id string) {
	if trailing != "/" {
		w.WriteHeader(http.StatusNotImplemented)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}

func (h *HivemindServer) apiSensorPut(w http.ResponseWriter, trailing, id string, body []byte) {
	value, _ := strconv.Atoi(string(body))
	err := h.store.storeSensorValue(id, value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusAccepted)
}
