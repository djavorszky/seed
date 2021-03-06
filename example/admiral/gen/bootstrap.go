package gen

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Service is the struct that will be exposed to serve HTTP traffic.
type Service struct {
	router      *mux.Router
	serviceImpl AdmiralService
}

// ServeHTTP is what ultimately allows this service to be used by the standard library's
// listen and serve functions
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Route is a struct that holds the path, the handler to be called when that path is hit, and
// which list of methods it should serve.
type route struct {
	path    string
	handler http.HandlerFunc
	methods []string
}

// New returns a new service implementation, using the service as a dependency. It also sets up the routes
// and the middlewares.
func New(service AdmiralService) *Service {
	s := &Service{
		router:      mux.NewRouter(),
		serviceImpl: service,
	}

	s.routes()
	s.middlewares()

	return s
}

// routes sets up the routes to be served by the service
func (s *Service) routes() {
	routes := []route{
		{
			path:    "/",
			handler: s.serviceImpl.Index(),
			methods: []string{http.MethodGet},
		},
	}

	for _, route := range routes {
		s.router.StrictSlash(true).
			HandleFunc(route.path, route.handler).
			Methods(route.methods...)
	}
}

// middlewares sets up the middlewares to be set up by the service
func (s *Service) middlewares() {
	mws := []mux.MiddlewareFunc{
		s.serviceImpl.LoggerMw,
	}

	for _, mw := range mws {
		s.router.Use(mw)
	}
}
