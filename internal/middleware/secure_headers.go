package middleware

import (
	"fmt"
	"net/http"
)

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
			"; script-src 'self' 'wasm-unsafe-eval' https://static.cloudflareinsights.com" +
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
