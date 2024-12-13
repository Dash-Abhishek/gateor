package main

import "net/http"

type ReqProcessorInterface interface {
	PreFlow(http.ResponseWriter, *http.Request)
	PostFlow(http.ResponseWriter, *http.Request)
}
