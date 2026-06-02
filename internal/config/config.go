// All settings are env-driven so the binary runs the same in any environment
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all runtime configuration
// Fields are populated from env vars and validated on Load
// The struct is read-only after Load returns
type Config struct {
	// HTTP server
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration

	// App
	Env         string // dev/prod
	SiteName    string
	BaseURL     string
	ContentDir  string
	TemplateDir string
	StaticDir   string

	// Security
	CookieSecret   []byte // 32+ bytes; used to sign/encrypt auth cookies
	CSPReportToURI string
	TrustedProxy   string // X-Forwarded-* respected only when peer matches

	// OIDC (all optional, when unset the app runs without auth)
	OIDCIssuer       string
	OIDCClientID     string
	OIDCClientSecret string
	OIDCRedirectURL  string
	OIDCScopes       []string
}

// Reads configuration from the process environment
func Load() (*Config, error) {
	c := &Config{
		Addr:            getenv("ADDR", ":8080"),
		ReadTimeout:     getenvDuration("READ_TIMEOUT", 10*time.Second),
		WriteTimeout:    getenvDuration("WRITE_TIMEOUT", 15*time.Second),
		ShutdownTimeout: getenvDuration("SHUTDOWN_TIMEOUT", 15*time.Second),

		Env:         getenv("APP_ENV", "dev"),
		SiteName:    getenv("SITE_NAME", "YOU FORGOT TO LOAD THE ENV FILE"),
		BaseURL:     getenv("BASE_URL", "http://localhost:8080"),
		ContentDir:  getenv("CONTENT_DIR", "content"),
		TemplateDir: getenv("TEMPLATE_DIR", "web/templates"),
		StaticDir:   getenv("STATIC_DIR", "web/static"),

		CSPReportToURI: getenv("CSP_REPORT_TO_URI", "/csp"),
		// CSPReportURI: getenv("CSP_REPORT_URI", ""),
		TrustedProxy: getenv("TRUSTED_PROXY", ""),

		OIDCIssuer:       getenv("OIDC_ISSUER", ""),
		OIDCClientID:     getenv("OIDC_CLIENT_ID", ""),
		OIDCClientSecret: getenv("OIDC_CLIENT_SECRET", ""),
		OIDCRedirectURL:  getenv("OIDC_REDIRECT_URL", ""),
		OIDCScopes:       splitCSV(getenv("OIDC_SCOPES", "openid,profile,email")),
	}

	secret := os.Getenv("COOKIE_SECRET")
	if secret == "" {
		// In dev we use a fake secret. In prod this must be set explicitly for consistancy
		if c.Env == "prod" {
			return nil, fmt.Errorf("COOKIE_SECRET is required in prod (32+ random bytes, base64 or hex)")
		}
		c.CookieSecret = []byte("dev-cookie-secret-not-for-production!!")
	} else {
		if len(secret) < 32 {
			return nil, fmt.Errorf("COOKIE_SECRET must be at least 32 bytes")
		}
		c.CookieSecret = []byte(secret)
	}

	return c, nil
}

// OIDCEnabled reports whether all required OIDC settings are present.
// When false, /protected falls back to a dev/disabled mode.
func (c *Config) OIDCEnabled() bool {
	return c.OIDCIssuer != "" && c.OIDCClientID != "" && c.OIDCClientSecret != "" && c.OIDCRedirectURL != ""
}

// IsProd reports whether the app is running in production mode.
func (c *Config) IsProd() bool { return c.Env == "prod" }

func getenv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func getenvDuration(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// Atoi parses an int env var with a fallback. Useful for callers that want
// numeric configuration without pulling in strconv directly.
func Atoi(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
