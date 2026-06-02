package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"ethanashley.net/website-main-go/internal/blog"
	"ethanashley.net/website-main-go/internal/config"
	"ethanashley.net/website-main-go/internal/handlers"
	"ethanashley.net/website-main-go/internal/middleware"
	"ethanashley.net/website-main-go/internal/templates"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("config load failed", "err", err)
		os.Exit(2)
	}

	if err := run(cfg, logger); err != nil {
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

func run(cfg *config.Config, logger *slog.Logger) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Templates
	tpl, err := templates.New(templates.Options{Dir: cfg.TemplateDir})
	if err != nil {
		return err
	}

	// Blog content
	blogStore := blog.New()
	if err := blogStore.Load(filepath.Join(cfg.ContentDir, "blog"), !cfg.IsProd()); err != nil {
		return err
	}
	logger.Info("blog loaded", "posts", len(blogStore.All()))

	// Auth
	auth, err := handlers.NewAuth(ctx, cfg, logger)
	if err != nil {
		return err
	}

	// Page handlers
	pages := &handlers.Pages{
		Cfg:   cfg,
		Tpl:   tpl,
		Blog:  blogStore,
		Auth:  auth,
		Log:   logger,
		Start: time.Now(),
	}

	mux := http.NewServeMux()
	// Static assets: served directly
	fileServer := http.FileServer(http.Dir(cfg.StaticDir))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	// PAGES
	mux.HandleFunc("GET /", pages.Home)
	mux.HandleFunc("GET /about", pages.About)
	mux.HandleFunc("GET /blog", pages.BlogIndex)
	mux.HandleFunc("GET /blog/{slug}", pages.BlogPost)
	mux.HandleFunc("GET /tools", pages.Tools)
	mux.HandleFunc("GET /protected", auth.RequireAuth(pages.Protected))

	// Feeds & SEO
	mux.HandleFunc("GET /feed.xml", pages.Feed)
	mux.HandleFunc("GET /atom.xml", pages.Feed)
	mux.HandleFunc("GET /robots.txt", pages.Robots)
	mux.HandleFunc("GET /sitemap.xml", pages.Sitemap)

	// OIDC
	mux.HandleFunc("GET /auth/login", auth.LoginHandler)
	mux.HandleFunc("GET /auth/callback", auth.CallbackHandler)
	mux.HandleFunc("POST /auth/logout", auth.LogoutHandler)
	mux.HandleFunc("GET /auth/logout", auth.LogoutHandler)

	// Security
	mux.HandleFunc("POST /csp", pages.CSPReport)

	// Health
	mux.HandleFunc("GET /healthz", pages.Healthz)
	mux.HandleFunc("GET /readyz", pages.Readyz)

	handler := middleware.Chain(mux,
		middleware.RequestID,
		middleware.Recover(logger),
		middleware.Logger(logger),
		middleware.SecureHeaders(middleware.SecureHeadersOptions{BaseURL: cfg.BaseURL, CSPReportToURI: cfg.CSPReportToURI, HSTS: cfg.IsProd()}),
		middleware.NoIndex("/protected", "/auth/"),
	)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 18, // 256 KiB
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("listening", "addr", cfg.Addr, "env", cfg.Env, "oidc", cfg.OIDCEnabled())
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-errCh:
		return err
	}

	shutCtx, shutCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutCancel()
	return srv.Shutdown(shutCtx)
}
