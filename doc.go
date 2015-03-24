/*
Package testmux provides an HTTP request multiplexer that can be used when
writing tests where is is necessary to validate that requests are received in
a specific order.

Any number of routes can be registered with the testmux.Router. Each route is
described by an HTTP method and URL path along with the handler that should be
invoked for a matching request.

Every time that a new request is dipatched via the ServeHTTP method, the router
will track whether the request matches one of the registered handlers and if it
was received in the correct sequence. The AssertVisited method can be used to
check for any 1) requests received out of order, 2) requests received for
unregistered routes, and 3) routes registered but never requested.

		func TestRequests(t *testing.T) {
			// Instantiate testmux router and register new routes
			router := &testmux.Router{}
			router.RegisterResp("GET", "/foo", 200, "Hello")
			router.RegisterResp("GET", "/bar", 200, "Bonjour")

			// Set-up test server
			server := httptest.NewServer(router)

			// Make requests
			http.Get(server.URL+"/foo")
			http.Get(server.URL+"/bar")

			// Assert that all registered routes visited in the correct order
			router.AssertVisited(t)
		}

When registering routes, user's can use either the RegisterResp method
which accepts a static response code and body to be returned for a
matching request, or the RegisterFunc method which takes a traditional
handler function.
*/
package testmux
