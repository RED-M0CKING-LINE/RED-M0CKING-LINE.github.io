// Implements optional OIDC login with stateless sessions stored in an encrypted cookie
package handlers

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"ethanashley.net/website-main-go/internal/config"
)

const (
	sessionCookie = "ea-web_session"
	stateCookie   = "ea-web_oidc_state"
	sessionTTL    = 12 * time.Hour
)

// Subset of OIDC claims we persist in the cookie
type Session struct {
	Sub     string    `json:"sub"`
	Email   string    `json:"email,omitempty"`
	Name    string    `json:"name,omitempty"`
	Expires time.Time `json:"exp"`
}

// Reports whether the session has not yet expired
func (s Session) Authed() bool { return !s.Expires.IsZero() && time.Now().Before(s.Expires) }

// Wires OIDC login/callback/logout and provides a RequireAuth middleware
// When OIDC is not configured it operates in "dev disabled" mode: the protected page renders an informational notice instead of redirecting
type Auth struct {
	cfg  *config.Config
	log  *slog.Logger
	key  []byte // 32 byte derived key for AES-GCM
	hmac []byte // 32 byte HMAC key for state cookie

	// Only populated when OIDCEnabled()
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	oauth    *oauth2.Config
}

// Initializes Auth
// If OIDC is configured, it dials the discovery endpoint up-front so misconfiguration fails at startup
func NewAuth(ctx context.Context, cfg *config.Config, log *slog.Logger) (*Auth, error) {
	sum := sha256.Sum256(append([]byte("enc:"), cfg.CookieSecret...))
	hsum := sha256.Sum256(append([]byte("mac:"), cfg.CookieSecret...))
	a := &Auth{cfg: cfg, log: log, key: sum[:], hmac: hsum[:]}

	if !cfg.OIDCEnabled() {
		log.Info("OIDC disabled, no provider configured,  /protected runs in DEV mode")
		return a, nil
	}

	provider, err := oidc.NewProvider(ctx, cfg.OIDCIssuer)
	if err != nil {
		return nil, fmt.Errorf("OIDC discover: %w", err)
	}
	a.provider = provider
	a.verifier = provider.Verifier(&oidc.Config{ClientID: cfg.OIDCClientID})
	a.oauth = &oauth2.Config{
		ClientID:     cfg.OIDCClientID,
		ClientSecret: cfg.OIDCClientSecret,
		RedirectURL:  cfg.OIDCRedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       cfg.OIDCScopes,
	}
	log.Info("OIDC configured", "issuer", cfg.OIDCIssuer)
	return a, nil
}

// Reports whether real OIDC flow is active
func (a *Auth) Enabled() bool { return a.oauth != nil }

// Starts the OIDC flow by setting a state cookie and redirecting to the provider PKCE is used to bind the code to this client
func (a *Auth) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if !a.Enabled() {
		http.Error(w, "OIDC disabled", http.StatusServiceUnavailable)
		return
	}
	state, err := randomToken(24)
	if err != nil {
		http.Error(w, "state", http.StatusInternalServerError)
		return
	}
	nonce, err := randomToken(24)
	if err != nil {
		http.Error(w, "nonce", http.StatusInternalServerError)
		return
	}

	// State cookie: state + nonce, HMAC-signed. CSRF protection
	payload := state + "|" + nonce
	mac := hmac.New(sha256.New, a.hmac)
	mac.Write([]byte(payload))
	signed := payload + "|" + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	http.SetCookie(w, &http.Cookie{
		Name:     stateCookie,
		Value:    signed,
		Path:     "/",
		HttpOnly: true,
		Secure:   a.cfg.IsProd(),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   600,
	})

	url := a.oauth.AuthCodeURL(state, oidc.Nonce(nonce))
	http.Redirect(w, r, url, http.StatusFound)
}

// Validates state, exchanges the code, verifies the ID token, extracts claims, and writes the session cookie
func (a *Auth) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	if !a.Enabled() {
		http.Error(w, "oidc disabled", http.StatusServiceUnavailable)
		return
	}
	ctx := r.Context()

	c, err := r.Cookie(stateCookie)
	if err != nil {
		http.Error(w, "missing state cookie", http.StatusBadRequest)
		return
	}
	// Clear state cookie immediately.
	http.SetCookie(w, &http.Cookie{Name: stateCookie, Value: "", Path: "/", MaxAge: -1})

	parts := splitN(c.Value, "|", 3)
	if len(parts) != 3 {
		http.Error(w, "bad state cookie", http.StatusBadRequest)
		return
	}
	state, nonce, sig := parts[0], parts[1], parts[2]
	mac := hmac.New(sha256.New, a.hmac)
	mac.Write([]byte(state + "|" + nonce))
	want := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(sig), []byte(want)) {
		http.Error(w, "state signature mismatch", http.StatusBadRequest)
		return
	}
	if got := r.URL.Query().Get("state"); got != state {
		http.Error(w, "state mismatch", http.StatusBadRequest)
		return
	}

	tok, err := a.oauth.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "code exchange failed", http.StatusBadRequest)
		return
	}
	rawID, ok := tok.Extra("id_token").(string)
	if !ok {
		http.Error(w, "no id_token", http.StatusBadRequest)
		return
	}
	idt, err := a.verifier.Verify(ctx, rawID)
	if err != nil {
		http.Error(w, "id_token verify", http.StatusUnauthorized)
		return
	}
	if idt.Nonce != nonce {
		http.Error(w, "nonce mismatch", http.StatusUnauthorized)
		return
	}

	var claims struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	_ = idt.Claims(&claims)

	sess := Session{
		Sub:     idt.Subject,
		Email:   claims.Email,
		Name:    claims.Name,
		Expires: time.Now().Add(sessionTTL),
	}
	if err := a.writeSession(w, sess); err != nil {
		http.Error(w, "session", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/protected", http.StatusFound)
}

// Clears session cookie and redirects home
func (a *Auth) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: sessionCookie, Value: "", Path: "/", MaxAge: -1, HttpOnly: true})
	http.Redirect(w, r, "/", http.StatusFound)
}

// SessionFrom reads and validates the session cookie, returning the decoded Session and whether it is present+valid
func (a *Auth) SessionFrom(r *http.Request) (Session, bool) {
	var s Session
	c, err := r.Cookie(sessionCookie)
	if err != nil {
		return s, false
	}
	raw, err := a.decrypt(c.Value)
	if err != nil {
		return s, false
	}
	if err := json.Unmarshal(raw, &s); err != nil {
		return s, false
	}
	return s, s.Authed()
}

// Gates a handler behind a valid session
// When OIDC is disabled the dev mode allows access but exposes a flag so templates can show a notice
func (a *Auth) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.Enabled() {
			next(w, r) // dev fallback and shows a "OIDC disabled" notice
			return
		}
		if _, ok := a.SessionFrom(r); !ok {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}
		next(w, r)
	}
}

// Encrypts and writes the session cookie
func (a *Auth) writeSession(w http.ResponseWriter, s Session) error {
	raw, err := json.Marshal(s)
	if err != nil {
		return err
	}
	enc, err := a.encrypt(raw)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    enc,
		Path:     "/",
		HttpOnly: true,
		Secure:   a.cfg.IsProd(),
		SameSite: http.SameSiteLaxMode,
		Expires:  s.Expires,
	})
	return nil
}

// Uses AES-GCM with a random nonce; nonce is prepended to ciphertext
func (a *Auth) encrypt(plain []byte) (string, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	ct := gcm.Seal(nonce, nonce, plain, nil)
	return base64.RawURLEncoding.EncodeToString(ct), nil
}

func (a *Auth) decrypt(s string) ([]byte, error) {
	data, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(data) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ct := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ct, nil)
}

func randomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// Split strings every N
func splitN(s, sep string, n int) []string {
	out := make([]string, 0, n)
	for i := 0; i < n-1; i++ {
		idx := indexOf(s, sep)
		if idx < 0 {
			break
		}
		out = append(out, s[:idx])
		s = s[idx+len(sep):]
	}
	out = append(out, s)
	return out
}

func indexOf(s, sep string) int {
	for i := 0; i+len(sep) <= len(s); i++ {
		if s[i:i+len(sep)] == sep {
			return i
		}
	}
	return -1
}
