package main

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// home returns a 404 placeholder page for non-existent paths.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.notFound(w)
}

// healthz is a dummy healthcheck
func (app *application) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (app *application) NewUpstream(name string, rawurl string) *upstream {
	target, err := url.Parse(rawurl)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	// Based on: https://golang.org/src/net/http/httputil/reverseproxy.go
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		// Replace the target host in case we want to use a virtual domain
		req.Host = target.Host
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		// Since we're targetting exactly one page, it's better to replace the path
		req.URL.Path = target.Path
		req.URL.RawPath = target.RawPath
		// Intentionally prohibit any changes in RawQuery
		req.URL.RawQuery = targetQuery
		if _, ok := req.Header["User-Agent"]; !ok {
			// Explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}

	// The values below are equal to http.DefaultTransport except for TLSClientConfig
	/* #nosec G402 */
	var transport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		ExpectContinueTimeout: 1 * time.Second,
	}

	proxy := &httputil.ReverseProxy{Director: director, Transport: transport}
	proxy.ErrorLog = app.errorLog
	proxy.FlushInterval = time.Millisecond * 200

	return &upstream{
		name:  name,
		url:   target,
		proxy: proxy,
	}
}
