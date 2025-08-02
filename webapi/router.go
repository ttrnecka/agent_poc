package main

import (
	"net/http"
	"strings"
)

type Middleware func(http.Handler) http.Handler

// Utility: compose multiple middleware
func chainMiddleware(h http.Handler, mws ...Middleware) http.Handler {
	for _, mw := range mws {
		h = mw(h)
	}
	return h
}

// Smart router with prefix-specific middleware
type MiddlewareRouter struct {
	mux        *http.ServeMux
	middleware map[string][]Middleware // path prefix -> middleware chain
	defaultMW  []Middleware
}

func NewMiddlewareRouter() *MiddlewareRouter {
	return &MiddlewareRouter{
		mux:        http.NewServeMux(),
		middleware: make(map[string][]Middleware),
	}
}

func (mr *MiddlewareRouter) Handle(pattern string, handler http.Handler) {
	mr.mux.Handle(pattern, handler)
}

func (mr *MiddlewareRouter) UseDefault(mw ...Middleware) {
	mr.defaultMW = mw
}

func (mr *MiddlewareRouter) UseForPrefix(prefix string, mw ...Middleware) {
	mr.middleware[prefix] = mw
}

func (mr *MiddlewareRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := mr.mux

	// Match longest prefix
	var chain []Middleware
	longest := 0
	for prefix, mws := range mr.middleware {
		if strings.HasPrefix(r.URL.Path, prefix) && len(prefix) > longest {
			chain = mws
			longest = len(prefix)
		}
	}

	if len(chain) == 0 {
		chain = mr.defaultMW
	}

	chainMiddleware(h, chain...).ServeHTTP(w, r)
}
