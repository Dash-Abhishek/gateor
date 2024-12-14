package internal

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func preFlow(w http.ResponseWriter, r *http.Request) {
	fmt.Println("executing default preflow")
}

func postFlow(q http.ResponseWriter, r *http.Request) {
	fmt.Println("executing default postflow")
}

func (svc service) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	preFlow(w, r)
	defer postFlow(w, r)

	targetUrl, err := url.Parse(svc.Target.Host)
	if err != nil {
		http.Error(w, "Invalid target url", http.StatusBadGateway)
		return
	}

	requestPath := r.URL.Path
	// strip basePath
	if svc.StripBasePath {
		requestPath = strings.TrimPrefix(r.URL.Path, svc.Path)
	}

	// TODO add rate limit

	// reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	proxy.Director = func(req *http.Request) {

		req.Header = r.Header.Clone() // Forward all headers
		req.Host = targetUrl.Host     // Set the Host header to the target host
		req.URL.Scheme = targetUrl.Scheme
		req.URL.Host = targetUrl.Host
		req.URL.Path = requestPath
	}

	proxy.ServeHTTP(w, r)

}
