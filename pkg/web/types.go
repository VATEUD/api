package web

import (
	"api/pkg/api/division"
	"api/pkg/api/news"
	"api/pkg/api/solo_phases"
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
	Path     string
	Methods  []string
	Function http.HandlerFunc
	Permission
}

type Middleware struct {
	Name     string
	Function mux.MiddlewareFunc
}

type Permission struct {
	AuthNeeded bool
	GuestOnly  bool
	AllowCors  bool
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
			Permission{
				false,
				true,
				false,
			},
		},
		{
			"/auth/validate",
			[]string{
				"GET",
			},
			connect.Validate,
			Permission{
				false,
				true,
				false,
			},
		},
		{
			"/api/user",
			[]string{
				"GET",
			},
			oauth2.User,
			Permission{
				true,
				false,
				true,
			},
		},
		{
			"/api/division/examiners",
			[]string{
				"GET",
			},
			division.Examiners,
			Permission{
				false,
				false,
				true,
			},
		},
		{
			"/api/division/instructors",
			[]string{
				"GET",
			},
			division.Instructors,
			Permission{
				false,
				false,
				true,
			},
		},
		{
			"/api/news",
			[]string{
				"GET",
			},
			news.NewsIndex,
			Permission{
				false,
				false,
				true,
			},
		},
		{
			"/api/news/{id}",
			[]string{
				"GET",
			},
			news.NewsShow,
			Permission{
				false,
				false,
				true,
			},
		},
		{
			"/api/subdivisions",
			[]string{
				"GET",
			},
			subdivision.Subdivisions,
			Permission{
				false,
				false,
				true,
			},
		},
		{
			"/api/subdivisions/view",
			[]string{
				"GET",
			},
			subdivision.Subdivisions,
			Permission{
				false,
				false,
				true,
			},
		},
		{
			"/api/subdivisions/view/{subdivision}",
			[]string{
				"GET",
			},
			subdivision.Subdivision,
			Permission{
				false,
				false,
				true,
			},
		},
		{
			"/api/subdivisions/instructors",
			[]string{
				"GET",
			},
			subdivision.Instructors,
			Permission{
				false,
				false,
				true,
			},
		},
		{
			"/api/subdivisions/instructors/{subdivision}",
			[]string{
				"GET",
			},
			subdivision.InstructorsFilter,
			Permission{
				false,
				false,
				true,
			},
		},
		{
			"/api/staff",
			[]string{
				"GET",
			},
			division.Staff,
			Permission{
				false,
				false,
				true,
			},
		},
		{
			"/api/events/view",
			[]string{
				"GET",
			},
			myvatsim.AllEvents,
			Permission{
				false,
				false,
				true,
			},
		},
		{
			"/api/events/view/{amount}",
			[]string{
				"GET",
			},
			myvatsim.EventsByAmount,
			Permission{
				false,
				false,
				true,
			},
		},
		{
			"/api/events/filter/days/{days}",
			[]string{
				"GET",
			},
			myvatsim.EventsFilterDays,
			Permission{
				false,
				false,
				true,
			},
		},
		{
			"/api/uploads/view",
			[]string{
				"GET",
			},
			uploads.List,
			Permission{
				false,
				false,
				true,
			},
		},
		{
			"/api/uploads/download/{id}",
			[]string{
				"GET",
			},
			uploads.Download,
			Permission{
				false,
				false,
				true,
			},
		},
		{
			"/api/uploads/filter/{type}",
			[]string{
				"GET",
			},
			uploads.Filter,
			Permission{
				false,
				false,
				true,
			},
		},
		{
			"/api/solo_phases",
			[]string{
				"GET",
			},
			solo_phases.RetrieveAll,
			Permission{
				false,
				false,
				true,
			},
		},
		{
			"/api/solo_phases/view",
			[]string{
				"GET",
			},
			solo_phases.RetrieveAll,
			Permission{
				false,
				false,
				true,
			},
		},
		{
			"/api/solo_phases/view/{subdivision}",
			[]string{
				"GET",
			},
			solo_phases.RetrieveBySubdivision,
			Permission{
				false,
				false,
				true,
			},
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
