package main

import "net/http"

type PluginInterface interface {
	Handle(w http.ResponseWriter, r *http.Request)
}
