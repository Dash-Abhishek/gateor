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

func preFlow(w http.ResponseWriter, r *http.Request) {
	pkg.Log.Debug("executing default preflow")
}

func postFlow(q http.ResponseWriter, r *http.Request, metaData reqMetaData) {
	status := http.StatusOK
	if w, ok := q.(interface{ Status() int }); ok {
		status = w.Status()
	}
	pkg.Log.Info("request", slog.String("Path", r.URL.Path), slog.Int("status", status), slog.String("Target", metaData.Target), slog.Any("Duration", metaData.TotalDuration.Milliseconds()), slog.Any("ProxyDuration", metaData.ProxyDuration.Milliseconds()))
}

func (svc service) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	requestStartTime := time.Now()

	preFlow(w, r)

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

	//  rate limit
	svc.RateLimiter.Handle(w, r)

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
	proxy.ServeHTTP(w, r)

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

	defer postFlow(w, r, metaData)

}
