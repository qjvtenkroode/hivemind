package main

import (
	"fmt"
	"net/http"
)

// HivemindServer is a HTTP interface for Hivemind
type HivemindServer struct {
	http.Handler
}

// NewHivemindServer creates a HivemindServer with routing configured
func NewHivemindServer() *HivemindServer {
	h := new(HivemindServer)

	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(h.rootHandler))
	router.Handle("/api/", http.HandlerFunc(h.apiHandler))
	router.Handle("/api/sensor/", http.HandlerFunc(h.apiSensorHandler))

	h.Handler = router

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
	w.Header().Set("content-type", "application/json")
	fmt.Fprint(w, "64")
}
