package main

import (
	"fmt"
	"net/url"
	"os"

	"gopkg.in/yaml.v2"
)

// isURL checks whether a string contains a valid URL
func (app *application) isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// loadEndpointsConfig returns a map with metrics endpoints to forward requests to
func (app *application) loadEndpointsConfig() (map[string]string, error) {
	rawConfig := make(map[string]string)
	config := make(map[string]string)

	yamlFile, err := os.ReadFile(app.ConfigPath)
	if err != nil {
		return rawConfig, err
	}

	err = yaml.Unmarshal(yamlFile, &rawConfig)
	if err != nil {
		return rawConfig, err
	}

	skipped := false
	for name, rawurl := range rawConfig {
		if !app.isURL(rawurl) {
			app.errorLog.Printf("skipped endpoint definition for %q, because %q is not a valid URL", name, rawurl)
			skipped = true
			continue
		}
		config[name] = rawurl
		app.infoLog.Printf("loaded endpoint definition %q: %q", name, rawurl)
	}

	if skipped {
		return config, fmt.Errorf("some of the endpoints were not loaded, so, please, check your configuration")
	}

	return config, nil
}
