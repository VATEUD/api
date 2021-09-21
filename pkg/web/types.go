package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	*http.Server
	handlers    []Handler
	router      *mux.Router
	middlewares []Middleware
}

type Handler struct {
	Path       string
	Methods    []string
	Function   http.HandlerFunc
	AuthNeeded bool
}

type Middleware struct {
	Name     string
	Function mux.MiddlewareFunc
}

func (server *Server) Start() error {
	server.registerMiddlewares()
	server.registerRoutes()
	return server.ListenAndServe()
}

func (server *Server) registerRoutes() {
	server.loadRoutes()
	for _, h := range server.handlers {
		server.router.HandleFunc(h.Path, h.Function).Methods(h.Methods...)
	}
	server.updateServerHandler()
}

func (server *Server) updateServerHandler() {
	server.Handler = server.router
}

func (server *Server) loadRoutes() {
	server.handlers = []Handler{
		{
			"/test",
			[]string{
				"GET",
			},
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("test"))
			},
			true,
		},
	}
}

func (server *Server) registerMiddlewares() {
	server.loadMiddlewares()
	for _, m := range server.middlewares {
		server.router.Use(m.Function)
	}
}

func (server *Server) loadMiddlewares() {
	server.middlewares = []Middleware{
		{
			Name:     "Rate limiting middleware",
			Function: rateLimitingMiddleware,
		},
		{
			Name:     "Authentication Middleware",
			Function: authMiddleware,
		},
	}
}

func (server Server) NeedsAuth(uri string) bool {
	for _, route := range server.handlers {
		if route.Path == uri {
			if route.AuthNeeded {
				return true
			}
			break
		}
	}

	return false
}
