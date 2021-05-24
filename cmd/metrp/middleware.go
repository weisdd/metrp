package main

import (
	"fmt"
	"net/http"
)

// allowedMethodsMiddleware allows only GET requests
func (app *application) allowedMethodsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			app.clientError(w, http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// basicAuthMiddleware verifies basic authentication
func (app *application) basicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.BasicAuth {
			next.ServeHTTP(w, r)
			return
		}

		username, password, ok := r.BasicAuth()
		if !ok || username != app.Username || password != app.Password {
			app.unathorized(w)
			return
		}

		// Do not forward login and password to upstream
		r.Header.Del("Authorization")

		next.ServeHTTP(w, r)
	})
}

// addBearerToken adds k8s token to all proxied requests
func (app *application) addBearerToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.UseToken {
			next.ServeHTTP(w, r)
			return
		}

		token := fmt.Sprintf("Bearer %s", app.Token)
		r.Header.Set("Authorization", token)

		next.ServeHTTP(w, r)
	})
}
