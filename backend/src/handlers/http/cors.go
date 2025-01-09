package http

import (
	"net/http"

	httpserver "github.com/w-h-a/pkg/serverv2/http"
)

var (
	// TODO: should come from config
	CORS = map[string]bool{
		"http://localhost:3000": true,
	}
)

type CORSMiddleware struct {
	handler http.Handler
}

func (m CORSMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	if CORS[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else if len(origin) > 0 && CORS["*"] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		return
	}

	m.handler.ServeHTTP(w, r)
}

func NewCORSMiddleware() httpserver.Middleware {
	return func(h http.Handler) http.Handler {
		return CORSMiddleware{h}
	}
}
