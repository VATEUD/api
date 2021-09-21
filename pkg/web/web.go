package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

const (
	defaultPortFieldName = "HTTP_PORT"
	defaultPort          = ":3000"
)

var server *Server

// New returns a new server structure
func New() *Server {
	server = &Server{
		&http.Server{
			Addr:         retrieveHttpPort(),
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
		},
		nil,
		mux.NewRouter(),
		nil,
	}
	return server
}

// retrieveHttpPort retrieves the server port from the environment file. If one isn't provided, it'll use the default one.
func retrieveHttpPort() string {
	if port := os.Getenv(defaultPortFieldName); port != "" {
		return fmt.Sprintf(":%s", port)
	}

	return defaultPort
}
