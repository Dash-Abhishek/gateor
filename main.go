package main

import (
	"context"
	"gateor/internal"
	"gateor/pkg"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	mux := internal.InitializeMux()
	internal.LoadSvc()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		pkg.Log.Error("Error starting gateor server", "Error", err)
	}
	defer listener.Close()

	server := http.Server{
		Handler: mux,
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		pkg.Log.Info("Shutting down gateor")
		if err := server.Shutdown(context.Background()); err != nil {
			pkg.Log.Error("Error shutting down server", "Error", err)
		}
		pkg.Log.Info("gateor shutdown complete")
	}()

	pkg.Log.Info("Gateor starting on ", "Address", listener.Addr().String())
	// Start server
	if err := server.Serve(listener); err != nil {
		pkg.Log.Error("Shutting server", "Error", err)
	}

}
