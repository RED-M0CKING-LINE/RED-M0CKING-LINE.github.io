// Page handlers
package handlers

import (
	"log/slog"
	"net/http"
	"runtime"
	"time"

	"ethanashley.net/website-main-go/internal/blog"
	"ethanashley.net/website-main-go/internal/config"
	"ethanashley.net/website-main-go/internal/templates"
)

// Bundles page handlers with the dependencies they need
type Pages struct {
	Cfg   *config.Config
	Tpl   *templates.Engine
	Blog  *blog.Store
	Auth  *Auth
	Log   *slog.Logger
	Start time.Time
}

// Returns a map of common template data merged with extras
func (p *Pages) base(r *http.Request, page string, extra map[string]any) map[string]any {
	sess, ok := p.Auth.SessionFrom(r)
	m := map[string]any{
		"SiteName":    p.Cfg.SiteName,
		"BaseURL":     p.Cfg.BaseURL,
		"Page":        page,
		"OIDCEnabled": p.Auth.Enabled(),
		"Authed":      ok,
		"Session":     sess,
		"Year":        time.Now().Year(),
	}
	for k, v := range extra {
		m[k] = v
	}
	return m
}

func (p *Pages) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		p.NotFound(w, r)
		return
	}
	posts := p.Blog.All()
	if len(posts) > 3 {
		posts = posts[:3]
	}
	_ = p.Tpl.Render(w, "home", p.base(r, "home", map[string]any{
		"RecentPosts": posts,
	}))
}

func (p *Pages) About(w http.ResponseWriter, r *http.Request) {
	_ = p.Tpl.Render(w, "about", p.base(r, "about", map[string]any{
		"Uptime":    time.Since(p.Start).Round(time.Second).String(),
		"GoVersion": runtime.Version(),
		"NumGorout": runtime.NumGoroutine(),
		"NumCPU":    runtime.NumCPU(),
	}))
}

func (p *Pages) BlogIndex(w http.ResponseWriter, r *http.Request) {
	_ = p.Tpl.Render(w, "blog_index", p.base(r, "blog", map[string]any{
		"Posts": p.Blog.All(),
	}))
}

func (p *Pages) BlogPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	post := p.Blog.Get(slug)
	if post == nil {
		p.NotFound(w, r)
		return
	}
	_ = p.Tpl.Render(w, "blog_post", p.base(r, "blog", map[string]any{
		"Post": post,
	}))
}

func (p *Pages) Tools(w http.ResponseWriter, r *http.Request) {
	_ = p.Tpl.Render(w, "tools", p.base(r, "tools", nil))
}

func (p *Pages) Protected(w http.ResponseWriter, r *http.Request) {
	sess, _ := p.Auth.SessionFrom(r)
	_ = p.Tpl.Render(w, "protected", p.base(r, "protected", map[string]any{
		"DevMode": !p.Auth.Enabled(),
		"User":    sess,
	}))
}

func (p *Pages) Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func (p *Pages) Readyz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ready"}`))
}

// Serves an Atom feed of all blog posts
func (p *Pages) Feed(w http.ResponseWriter, r *http.Request) {
	xml, err := p.Blog.AtomFeed(blog.FeedOptions{
		SiteName: p.Cfg.SiteName,
		BaseURL:  p.Cfg.BaseURL,
	})
	if err != nil {
		p.Log.Error("feed render", "err", err)
		http.Error(w, "feed error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/atom+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.Write(xml)
}

// Robots and Sitemap are intentionally minimal
// Nginx also serves these
func (p *Pages) Robots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("User-agent: *\nDisallow: /protected\nDisallow: /auth/\nSitemap: " + p.Cfg.BaseURL + "/sitemap.xml\n"))
}

func (p *Pages) Sitemap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	base := p.Cfg.BaseURL
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>`))
	w.Write([]byte(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`))
	for _, u := range []string{"/", "/blog", "/tools", "/about"} {
		w.Write([]byte("<url><loc>" + base + u + "</loc></url>"))
	}
	for _, post := range p.Blog.All() {
		w.Write([]byte("<url><loc>" + base + "/blog/" + post.Slug + "</loc><lastmod>" + post.Meta.Date.UTC().Format("2006-01-02") + "</lastmod></url>"))
	}
	w.Write([]byte(`</urlset>`))
}

func (p *Pages) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_ = p.Tpl.Render(w, "404", p.base(r, "404", nil))
}
