package testmux // import "github.com/CenturyLinkLabs/testmux"

import (
	"fmt"
	"net/http"
	"testing"
)

// Router registers routes to be matched and dispatches the associated
// handler. It will track handled requests to ensure that they are received
// in the same order in which they were originally registered.
//
// It implements the http.Hander interface.
type Router struct {
	routes []route
	index  int
	errors []string
}

type route struct {
	method  string
	path    string
	handler func(http.ResponseWriter, *http.Request)
	visited bool
}

// RegisterFunc registers a handler function for the given request method and
// path.
func (r *Router) RegisterFunc(method, path string, handler func(http.ResponseWriter, *http.Request)) {
	rte := route{method: method, path: path, handler: handler}
	r.routes = append(r.routes, rte)
}

// RegisterResp registers a static status code and body string to be returned
// for the given request method and path.
func (r *Router) RegisterResp(method, path string, status int, body string) {
	rte := route{
		method: method,
		path:   path,
		handler: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
			fmt.Fprintln(w, body)
		},
	}

	r.routes = append(r.routes, rte)
}

// ServeHTTP dispatches the handler registered in the matched route and tracks
// whether the route was requested in the correct order.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rte, ordered := r.match(req.Method, req.URL.Path)

	if rte != nil {
		rte.execute(w, req)

		if !ordered {
			r.addErrorf("Request out of order: %s %s", req.Method, req.URL.Path)
		}
	} else {
		http.NotFound(w, req)
		r.addErrorf("Unexpected request: %s %s", req.Method, req.URL.Path)
	}

	r.index++
}

// AssertVisited asserts that all of the registered routes were visited in the
// correct order. Conditions that will cause the assertion to fail include:
// requests received out of order, requests received for unregistered routes,
// and routes registered but never requested.
//
// Returns whether the assertion was successful (true) or not (false).
func (r *Router) AssertVisited(t *testing.T) bool {
	for _, rte := range r.routes {
		if !rte.visited {
			r.addErrorf("Unvisited route: %s %s", rte.method, rte.path)
		}
	}

	for _, err := range r.errors {
		t.Error(err)
	}

	return len(r.errors) == 0
}

// Adds an error string to the internal collection.
func (r *Router) addErrorf(format string, a ...interface{}) {
	r.errors = append(r.errors, fmt.Sprintf(format, a...))
}

// Given an HTTP method and request path, looks for a matching handler which
// has not already been visited. Returns the handler along with a flag
// indicating whether or not the handler is being invoked in the correct
// order.
func (r *Router) match(method, path string) (*route, bool) {
	for i, rte := range r.routes {
		if !rte.visited && rte.method == method && rte.path == path {
			return &r.routes[i], (i == r.index)
		}
	}

	return nil, false
}

// Execute the handler associated with the route and mark it as visited.
func (rte *route) execute(w http.ResponseWriter, req *http.Request) {
	rte.handler(w, req)
	rte.visited = true
}
