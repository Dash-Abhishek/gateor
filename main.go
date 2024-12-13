package main

import (
	"log"
	"net"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	trafficHandler := &TrafficHandler{mux: mux}
	services := LoadServices()

	for _, svc := range services {
		mux.Handle("/"+svc.Path, svc)
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	if err := http.Serve(listener, trafficHandler); err != nil {
		log.Fatalf("Error serving: %v", err)
	}

}

type TrafficHandler struct {
	mux *http.ServeMux
}

func (h *TrafficHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if h.mux != nil {
		h.mux.ServeHTTP(w, r)
		return
	}
	http.NotFound(w, r)

}
