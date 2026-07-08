// Package middleware provides HTTP middleware for the kyudo-dojo-hub API server.
package middleware

import (
	"net/http"
	"strings"
)

// DefaultCORSOrigin is used when CORS_ORIGIN is not set.
const DefaultCORSOrigin = "http://localhost:5173"

// ParseOrigins splits a comma-separated CORS_ORIGIN value into a list of
// trimmed, non-empty origins, falling back to the default single origin when
// the value is empty or blank.
func ParseOrigins(raw string) []string {
	origins := make([]string, 0, 1)
	for _, part := range strings.Split(raw, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	if len(origins) == 0 {
		return []string{DefaultCORSOrigin}
	}
	return origins
}

// CORS adds Cross-Origin Resource Sharing headers for the allowed origins.
//
// Multiple origins are supported: the request's Origin header is matched
// against the allow list and, when allowed, echoed back in
// Access-Control-Allow-Origin (a wildcard cannot be used together with
// credentials). Vary: Origin is set so caches key on the request origin.
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowed[o] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if _, ok := allowed[origin]; ok && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Add("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
