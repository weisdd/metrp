package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// clientError sends responses like 400 "Bad Request" to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound sends a 404 Not Found response to the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// notFound sends a 401 to the user.
func (app *application) unathorized(w http.ResponseWriter) {
	app.clientError(w, http.StatusUnauthorized)
}

// checkBasicAuthConfig checks whether all requirements for basic auth configuration are fullfilled
func (app *application) checkBasicAuthConfig() error {
	if !app.BasicAuth {
		app.infoLog.Print("authentication is disabled")
	} else {
		app.infoLog.Print("authentication is enabled")
		if app.Username == "" {
			return fmt.Errorf("please, specify a username (METRP_USERNAME)")
		}
		if app.Password == "" {
			return fmt.Errorf("please, specify a password (METRP_PASSWORD)")
		}
	}
	return nil
}

// getToken return k8s token, which can be later used for authorization
func (app *application) getToken() (string, error) {
	/* #nosec G101 */
	tokenPath := "/var/run/secrets/kubernetes.io/serviceaccount/token"
	tokenRaw, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", err
	}
	token := string(tokenRaw)
	token = strings.TrimSpace(token)
	return token, nil
}
