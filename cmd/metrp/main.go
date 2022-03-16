package main

import (
	"log"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

type application struct {
	errorLog                *log.Logger
	infoLog                 *log.Logger
	debugLog                *log.Logger
	upstreams               []*upstream
	IP                      string
	Port                    int           `env:"METRP_PORT" envDefault:"8080"`
	PreferredIPv4           string        `env:"METRP_PREFERRED_IPV4"`
	PreferredIPv4Prefix     string        `env:"METRP_PREFERRED_IPV4_PREFIX"`
	ConfigPath              string        `env:"METRP_CONFIG_PATH" envDefault:"metrp.yaml"`
	ReadTimeout             time.Duration `env:"METRP_READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout            time.Duration `env:"METRP_WRITE_TIMEOUT" envDefault:"10s"`
	GracefulShutdownTimeout time.Duration `env:"METRP_SHUTDOWN_TIMEOUT" envDefault:"20s"`
	BasicAuth               bool          `env:"METRP_BASIC_AUTH"`
	Username                string        `env:"METRP_USERNAME"`
	Password                string        `env:"METRP_PASSWORD"`
	UseToken                bool          `env:"METRP_USE_TOKEN" envDefault:"true"`
	Token                   string
}

type upstream struct {
	name  string
	url   *url.URL
	proxy *httputil.ReverseProxy
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	debugLog := log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		debugLog: debugLog,
	}

	err := env.Parse(app)
	if err != nil {
		app.errorLog.Fatalf("%+v\n", err)
	}

	err = app.checkBasicAuthConfig()
	if err != nil {
		app.errorLog.Fatal(err)
	}

	endpoints, err := app.loadEndpointsConfig()
	if err != nil {
		app.errorLog.Fatal(err)
	}

	for k, v := range endpoints {
		app.upstreams = append(app.upstreams, app.NewUpstream(k, v))
	}

	app.IP, err = app.getPreferredIP()
	if err != nil {
		app.errorLog.Fatal(err)
	}

	if app.UseToken {
		app.Token, err = app.getToken()
		if err != nil {
			app.errorLog.Fatal(err)
		}
	}

	err = app.serve()
	if err != nil {
		app.errorLog.Fatal(err)
	}
}
