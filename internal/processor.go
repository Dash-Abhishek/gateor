package internal

import (
	"gateor/pkg"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type reqMetaData struct {
	StartTime     time.Time
	ProxyDuration time.Duration
	TotalDuration time.Duration
	Target        string
}

type responseWriter struct {
	http.ResponseWriter
	written bool
	status  int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.written = true
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.written = true
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) Written() bool {
	return rw.written
}

func (svc service) preFlow(w http.ResponseWriter, r *http.Request) {
	svc.PluginChain.Handle(w, r)
}

func (svc service) postFlow(q http.ResponseWriter, r *http.Request, metaData reqMetaData) {
	status := http.StatusOK
	if w, ok := q.(interface{ Status() int }); ok {
		status = w.Status()
	}
	pkg.Log.Info("request", slog.String("Path", r.URL.Path), slog.Int("status", status), slog.String("Target", metaData.Target), slog.Any("Duration", metaData.TotalDuration.Milliseconds()), slog.Any("ProxyDuration", metaData.ProxyDuration.Milliseconds()))
}

func (svc service) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	requestStartTime := time.Now()

	respWriter := &responseWriter{ResponseWriter: w}

	svc.preFlow(respWriter, r)
	if respWriter.Written() {
		return
	}

	targetUrl, err := url.Parse(svc.Target.Host)
	if err != nil {
		http.Error(w, "Invalid target url", http.StatusBadGateway)
		return
	}

	// strip basePath
	requestPath := r.URL.Path
	if svc.StripBasePath {
		requestPath = strings.TrimPrefix(r.URL.Path, svc.Path)
	}

	// reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	proxy.Director = func(req *http.Request) {

		req.Header = r.Header.Clone() // Forward all headers
		req.Host = targetUrl.Host     // Set the Host header to the target host
		req.URL.Scheme = targetUrl.Scheme
		req.URL.Host = targetUrl.Host
		req.URL.Path = requestPath
	}

	proxyStartTime := time.Now()
	proxy.ServeHTTP(respWriter, r)

	// total time taken by target
	proxyDuration := time.Since(proxyStartTime)
	// turn around time
	totalDuration := time.Since(requestStartTime)

	metaData := reqMetaData{
		StartTime:     requestStartTime,
		ProxyDuration: proxyDuration,
		TotalDuration: totalDuration,
		Target:        svc.Target.Host,
	}

	go svc.postFlow(respWriter, r, metaData)

}
