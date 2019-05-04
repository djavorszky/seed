package gen

import "net/http"

// NameService encapsulates the handler interface, which holds all the methods to be called
// by the server, and middleware interface, which contains all the middlewares to be added to the service.
//
// "Name" here is a placeholder and should match the name of the module in the generated project
type NameService interface {
	NameHandler
	NameMiddleware
}

// NameHandler is the interface for the handlers. Any new endpoint added by seed will be added here as a
// new method on the interface.
//
// "Name" here is a placeholder and should match the name of the module in the generated project
type NameHandler interface {
	Index() http.HandlerFunc
}

// NameMiddleware is the interface for all the middlewares that will be added to all of the paths.
// May refactor to allow for adding to subpaths in the future.
//
// "Name" here is a placeholder and should match the name of the module in the generated project
type NameMiddleware interface {
	LoggerMw(http.Handler) http.Handler
}
