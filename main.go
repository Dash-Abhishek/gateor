package main

import (
	"gateor/internal"
	"log"
	"net"
	"net/http"
)

func main() {

	mux := internal.InitializeMux()
	internal.LoadSvc()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	if err := http.Serve(listener, mux); err != nil {
		log.Fatalf("Error serving: %v", err)
	}

}
