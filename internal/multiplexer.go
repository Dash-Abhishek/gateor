package internal

import (
	"net/http"
	"sync"
)

var muxOnce sync.Once
var mux *http.ServeMux

func InitializeMux() *http.ServeMux {

	muxOnce.Do(func() {
		mux = http.NewServeMux()
	})

	return mux

}
