package web

import (
	"api/pkg/api/division"
	"api/pkg/api/news"
	"api/pkg/api/subdivision"
	"api/pkg/api/uploads"
	"api/pkg/oauth2"
	"api/pkg/response"
	"api/pkg/vatsim/connect"
	"api/pkg/vatsim/myvatsim"
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
	AllowCors  bool
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
			"/auth/login",
			[]string{
				"GET",
			},
			connect.Login,
			false,
			true,
			false,
		},
		{
			"/auth/validate",
			[]string{
				"GET",
			},
			connect.Validate,
			false,
			true,
			false,
		},
		{
			"/api/user",
			[]string{
				"GET",
			},
			oauth2.User,
			true,
			false,
			true,
		},
		{
			"/api/division/examiners",
			[]string{
				"GET",
			},
			division.Examiners,
			false,
			false,
			true,
		},
		{
			"/api/division/instructors",
			[]string{
				"GET",
			},
			division.Instructors,
			false,
			false,
			true,
		},
		{
			"/api/news",
			[]string{
				"GET",
			},
			news.NewsIndex,
			false,
			false,
			true,
		},
		{
			"/api/news/{id}",
			[]string{
				"GET",
			},
			news.NewsShow,
			false,
			false,
			true,
		},
		{
			"/api/subdivisions",
			[]string{
				"GET",
			},
			subdivision.Subdivisions,
			false,
			false,
			true,
		},
		{
			"/api/subdivisions/view",
			[]string{
				"GET",
			},
			subdivision.Subdivisions,
			false,
			false,
			true,
		},
		{
			"/api/subdivisions/view/{subdivision}",
			[]string{
				"GET",
			},
			subdivision.Subdivision,
			false,
			false,
			true,
		},
		{
			"/api/subdivisions/instructors",
			[]string{
				"GET",
			},
			subdivision.Instructors,
			false,
			false,
			true,
		},
		{
			"/api/subdivisions/instructors/{subdivision}",
			[]string{
				"GET",
			},
			subdivision.InstructorsFilter,
			false,
			false,
			true,
		},
		{
			"/api/staff",
			[]string{
				"GET",
			},
			division.Staff,
			false,
			false,
			true,
		},
		{
			"/api/events/view",
			[]string{
				"GET",
			},
			myvatsim.AllEvents,
			false,
			false,
			true,
		},
		{
			"/api/events/view/{amount}",
			[]string{
				"GET",
			},
			myvatsim.EventsByAmount,
			false,
			false,
			true,
		},
		{
			"/api/events/filter/days/{days}",
			[]string{
				"GET",
			},
			myvatsim.EventsFilterDays,
			false,
			false,
			true,
		},
		{
			"/api/uploads/view",
			[]string{
				"GET",
			},
			uploads.List,
			false,
			false,
			true,
		},
		{
			"/api/uploads/download/{id}",
			[]string{
				"GET",
			},
			uploads.Download,
			false,
			false,
			true,
		},
		{
			"/api/uploads/filter/{type}",
			[]string{
				"GET",
			},
			uploads.Filter,
			false,
			false,
			true,
		},
		{
			"/app/authorize",
			[]string{
				"GET",
			},
			oauth2.Authorize,
			false,
			true,
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

func (server Server) AllowCors(uri string) bool {
	for _, route := range server.handlers {
		if route.Path == uri {
			if route.AllowCors {
				return true
			}
			break
		}
	}

	return false
}
