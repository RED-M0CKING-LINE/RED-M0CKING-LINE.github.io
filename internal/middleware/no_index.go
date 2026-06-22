package middleware

import (
	"net/http"
	"strings"
)

// Set X-Robots-Tag for paths that should never be indexed
// Nginx also sets this
func NoIndex(prefixes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, p := range prefixes {
				if strings.HasPrefix(r.URL.Path, p) {
					w.Header().Set("X-Robots-Tag", "noindex, nofollow")
					break
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
