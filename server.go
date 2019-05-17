package main

import (
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

	h.Handler = router

	return h
}

func (h *HivemindServer) rootHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.RequestURI()
	if url == "/" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (h *HivemindServer) apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	if r.URL.RequestURI() != "/api/" {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
