package seed

// ServiceDescriptor describes what the service should look like, and generates
// the output based on it.
type ServiceDescriptor struct {
	// Info holds generic information about the service itself.
	Info

	// Routes is a slice of route objects, which detail the endpoints on
	// which the service accepts requests, and to which service implementation
	// methods it should forward them.
	Routes []Route

	// Middlewares is a slice of middleware objects, which detail the
	// middlewares that should be added to the paths. The priority field
	// decides the order in which the middlewares are added. By default,
	// the order corresponds to the order the middlewares were added in.
	Middlwares []Middleware
}

// Route is an object that details an endpoint on which the service should serve
// content, what the implementation handler should be called, and which http
// methods are accepted.
type Route struct {
	// Info holds generic information about the route itself.
	Info

	// Path is the URI endpoint, relative to the root URL. Should start with a
	// slash. To add an endpoint that should listen on the root URL, specify
	// "/"
	Path string

	// StrictSlash decides what should happen when the path does not exactly
	// correspond to the request URI with regards to the trailing slash.
	// For example, having a route with "/path", accessing "/path/" will:
	// - StrictSlash false: 404 not found
	// - StrictSlash true: 301 moved permanently to "/path"
	StrictSlash bool

	// HttpMethods is a list of strings that the route supports. The contents
	// should correspond to the default HTTP methods: GET, POST, PUT, PATCH,
	// DELETE, OPTIONS, HEAD, CONNECT, and TRACE
	HttpMethods []string

	// HandlerName is the name of the method that will be called by the server
	// when the specified endpoint is requested, and as such, will contain the
	// business logic
	HandlerName string
}

type Middleware struct {
	Info
	Path        string
	HandlerName string
	Priority    int
}

// Info contains the generic information about a service object - its name, a
// short summary, and an optional longer description.
type Info struct {
	// Name should be a simple name for the service object.
	Name string

	// Summary should contain a short description of the service object,
	// typically a single sentence
	Summary string

	// Description may contain a longer explanation about what the service
	// object is for. It may provide more clarity to the reason for its
	// existence, or provide additional details about how it should be used.
	Description string
}
