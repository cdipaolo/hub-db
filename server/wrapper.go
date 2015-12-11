package server

import (
	"fmt"
	"net/http"
)

// HTTPHandlerWithMethod takes in a handler func and wraps
// it so only requests with the specified method (POST,
// GET, etc) will go through to the handler.
func HTTPHandlerWithMethod(method string, h http.HandlerFunc) http.HandlerFunc {
	return func(r http.Response, req http.Request) {
		r.Header().Add("Content-Type", "application/json")
		if req.Method != method {
			r.WriteHeader(http.StatusMethodNotAllowed)
			r.Write([]byte(fmt.Sprintf(`{"message": "Given method not allowed", "method": %v}`, req.Method)))
			return
		}

		h(r, req)
	}
}

/* Currently only need GET but this can be extended easily */

// Get takes in a HandlerFunc and makes it
// only accept GET requests
func Get(h http.HandlerFunc) http.HandlerFunc {
	return HTTPHandlerWithMethod("GET", h)
}
