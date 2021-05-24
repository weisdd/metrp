package main

import (
	"github.com/gorilla/mux"
)

// routes returns a router with all paths
func (app *application) routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", app.home)
	r.HandleFunc("/healthz", app.healthz)
	s := r.PathPrefix("/metrics").Subrouter()
	s.Use(app.allowedMethodsMiddleware)
	s.Use(app.basicAuthMiddleware)
	s.Use(app.addBearerToken)

	for _, upstream := range app.upstreams {
		s.HandleFunc("/"+upstream.name, upstream.proxy.ServeHTTP)
	}

	return r
}
