// Provides composable HTTP middleware: security headers, request logging, panic recovery, request ID, and a simple chain helper
package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

// Composes middlewares right-to-left so the call order matches the declared order (first declared = outermost)
func Chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// Catches panics and returns a 500
// The panic is logged with the stack so operators can diagnose without leaking details to clients
func Recover(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic",
						"err", rec,
						"path", r.URL.Path,
						"stack", string(debug.Stack()),
					)
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// Capture status code for logging
type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (s *statusRecorder) WriteHeader(c int) { s.status = c; s.ResponseWriter.WriteHeader(c) }
func (s *statusRecorder) Write(b []byte) (int, error) {
	if s.status == 0 {
		s.status = 200
	}
	n, err := s.ResponseWriter.Write(b)
	s.bytes += n
	return n, err
}

// Emits one structured log line per request
func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w}
			next.ServeHTTP(rec, r)
			logger.Info("http",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"bytes", rec.bytes,
				"dur_ms", time.Since(start).Milliseconds(),
				"rid", w.Header().Get("X-Request-ID"),
			)
		})
	}
}

// Assigns each request a short hex ID and echoes it in a header
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			var b [8]byte
			_, _ = rand.Read(b[:])
			id = hex.EncodeToString(b[:])
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}

// Configure the SecureHeaders middleware
type SecureHeadersOptions struct {
	// Content-Security-Policy header
	CSP string
	// CSP Report To header
	BaseURL        string
	CSPReportToURI string
	// Send Strict-Transport-Security when true (use behind TLS only)
	HSTS bool
}

// Adds response headers
// Nginx also impliments these, just to be sure
func SecureHeaders(opts SecureHeadersOptions) func(http.Handler) http.Handler {
	csp := opts.CSP
	if csp == "" {
		csp = "default-src 'self'" +
			"; img-src 'self' data:" +
			"; style-src 'self'" +
			"; script-src 'self' 'wasm-unsafe-eval' /cdn-cgi/" +
			// "; script-src 'self'" + //TESTING CSP
			"; font-src 'self'" +
			"; connect-src 'self'" +
			"; frame-ancestors 'none'" +
			"; base-uri 'self'" +
			"; form-action 'self'" +
			"; report-to csp-endpoint" +
			"; report-uri " + opts.CSPReportToURI
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			h.Set("Content-Security-Policy", csp)
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("X-Frame-Options", "DENY")
			h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
			h.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=()")
			h.Set("Cross-Origin-Opener-Policy", "same-origin")
			h.Set("Cross-Origin-Resource-Policy", "same-origin")
			if opts.HSTS {
				h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
			}
			if opts.BaseURL+opts.CSPReportToURI != "" {
				h.Set("Report-To", fmt.Sprintf(`[{"group":"csp-endpoint","max_age":86400,"endpoints":[{"url":"%s"}]}]`, opts.BaseURL+opts.CSPReportToURI))
				h.Set("Reporting-Endpoints", fmt.Sprintf(`csp-endpoint="%s"`, opts.BaseURL+opts.CSPReportToURI))
			}
			next.ServeHTTP(w, r)
		})
	}
}

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
