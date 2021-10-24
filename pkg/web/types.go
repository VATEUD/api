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
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
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
	AuthNeeded       bool
	GuestOnly        bool
	AllowCors        bool
	SubdivisionToken bool
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
			},
		},
		{
			"/api/solo_phases/create",
			[]string{
				"POST",
			},
			solo_phases.Create,
			Permission{
				true,
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
		uri = server.checkURI(route, uri)

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
		uri = server.checkURI(route, uri)

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
		uri = server.checkURI(route, uri)

		if route.Path == uri {
			if route.AllowCors {
				return true
			}
			break
		}
	}

	return false
}

func (server Server) NeedsSubdivisionToken(uri string) bool {
	for _, route := range server.handlers {
		uri = server.checkURI(route, uri)

		if route.Path == uri {
			if route.SubdivisionToken {
				return true
			}
			break
		}
	}

	return false
}

func (server Server) indexOfBracket(uri []string) []int {
	var indexes []int

	for i, data := range uri {
		if strings.Contains(data, "{") {
			indexes = append(indexes, i)
		}
	}

	return indexes
}

func (server Server) compileURI(parts []string) string {
	var uri string

	for _, part := range parts[1:] {
		uri += fmt.Sprintf("/%s", part)
	}

	return uri
}

func (server Server) checkURI(route Handler, uri string) string {
	if strings.Contains(route.Path, "{") {
		if strings.Contains(uri, route.Path[:strings.Index(route.Path, "{")]) {
			parts := strings.Split(uri, "/")
			pathParts := strings.Split(route.Path, "/")

			indexes := server.indexOfBracket(pathParts)

			for _, i := range indexes {
				parts[i] = pathParts[i]
			}

			return server.compileURI(parts)
		}
	}

	return uri
}
