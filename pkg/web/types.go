package web

import (
	"auth/pkg/api/division"
	"auth/pkg/oauth2"
	"auth/pkg/response"
	"auth/pkg/vatsim/connect"
	"github.com/gorilla/mux"
	"log"
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
	GuestOnly  bool
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
	log.Println("Registering the routes")
	server.loadRoutes()
	for _, h := range server.handlers {
		server.router.HandleFunc(h.Path, h.Function).Methods(h.Methods...)
	}
	server.router.NotFoundHandler = http.HandlerFunc(response.NotFoundHandler)
	server.router.MethodNotAllowedHandler = http.HandlerFunc(response.MethodNotAllowedHandler)
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
			false,
		},
		{
			"/auth/login",
			[]string{
				"GET",
			},
			connect.Login,
			false,
			true,
		},
		{
			"/auth/validate",
			[]string{
				"GET",
			},
			connect.Validate,
			false,
			true,
		},
		{
			"/api/user",
			[]string{
				"GET",
			},
			oauth2.User,
			true,
			false,
		},
		{
			"/api/division/examiners",
			[]string{
				"GET",
			},
			division.Examiners,
			false,
			false,
		},
	}
}

func (server *Server) registerMiddlewares() {
	log.Println("Registering the middlewares")
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

func (server Server) GuestOnly(uri string) bool {
	for _, route := range server.handlers {
		if route.Path == uri {
			if route.GuestOnly {
				return true
			}
			break
		}
	}

	return false
}
