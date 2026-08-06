package router

import (
	"net/http"
)

// Router to store all routes in map
type Router struct {
	routes map[string]http.HandlerFunc
}

// New() creates and returns a new Router instance
func New() *Router {
	return &Router{
		routes: make(map[string]http.HandlerFunc),
	}
}

// Handle method adds a route to the map when called
func (r *Router) Handle(path string, handler http.HandlerFunc) {
	r.routes[path] = handler
}

// ServeHTTP is called when a request comes in
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// Try to find the handler for this path
	handler, exists := r.routes[req.URL.Path]
	if !exists {
		http.NotFound(w, req) // Path not found send 404
		return
	}
	handler(w, req) // Run the handler when path is found
}
