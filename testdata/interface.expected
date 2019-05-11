package gen

import "net/http"

// AdmiralService encapsulates the handler interface, which holds all the methods to be called
// by the server, and middleware interface, which contains all the middlewares to be added to the service.
type AdmiralService interface {
	AdmiralHandler
	AdmiralMiddleware
}

// AdmiralHandler is the interface for the handlers. Any new endpoint added by seed will be added here as a
// new method on the interface.
type AdmiralHandler interface {
	Index() http.HandlerFunc
}

// AdmiralMiddleware is the interface for all the middlewares that will be added to all of the paths.
type AdmiralMiddleware interface {
	LoggerMw(http.Handler) http.Handler
}
