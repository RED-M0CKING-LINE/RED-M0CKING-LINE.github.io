// Loads, parses, and renders Markdown blog posts from the content directory
// Posts are loaded once on startup, so a restart is required for new content. This keeps it stateless
package blog

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

// YAML metadata parsed from the top of each post
type FrontMatter struct {
	Title   string    `yaml:"title"`
	Date    time.Time `yaml:"date"`
	Updated time.Time `yaml:"updated,omitempty"`
	Summary string    `yaml:"summary"`
	Tags    []string  `yaml:"tags"`
	Author  string    `yaml:"author"`
	Draft   bool      `yaml:"draft"`
}

// A rendered blog post
type Post struct {
	Slug    string
	Meta    FrontMatter
	HTML    template.HTML // sanitized HTML body
	Raw     string        // original markdown (without front matter metadata)
	ModTime time.Time     // file mtime
}

// Holds all parsed posts. It is read-only after Load
type Store struct {
	mu     sync.RWMutex
	posts  []*Post
	bySlug map[string]*Post
}

// Creates an empty store
func New() *Store {
	return &Store{bySlug: map[string]*Post{}}
}

// Load scans dir for *.md files, parses front matter, renders HTML, and populates the store
// Drafts are skipped unless includeDrafts is true
func (s *Store) Load(dir string, includeDrafts bool) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read blog dir: %w", err)
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.Typographer,
		),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithHardWraps(), html.WithXHTML()),
	)

	// Strict sanitizer: allow Markdown, block scripts/iframes/etc
	policy := bluemonday.UGCPolicy()
	policy.AllowAttrs("class").OnElements("code", "pre", "span", "div")
	policy.AllowAttrs("id").OnElements("h1", "h2", "h3", "h4", "h5", "h6")

	posts := make([]*Post, 0, len(entries))
	bySlug := make(map[string]*Post, len(entries))

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		fi, _ := e.Info()

		fm, body, err := parseFrontMatter(raw)
		if err != nil {
			return fmt.Errorf("front matter %s: %w", path, err)
		}
		if fm.Draft && !includeDrafts {
			continue
		}

		var buf bytes.Buffer
		if err := md.Convert(body, &buf); err != nil {
			return fmt.Errorf("render %s: %w", path, err)
		}
		clean := policy.SanitizeBytes(buf.Bytes())

		slug := strings.TrimSuffix(e.Name(), ".md")
		p := &Post{
			Slug:    slug,
			Meta:    fm,
			HTML:    template.HTML(clean),
			Raw:     string(body),
			ModTime: fi.ModTime(),
		}
		posts = append(posts, p)
		bySlug[slug] = p
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Meta.Date.After(posts[j].Meta.Date)
	})

	s.mu.Lock()
	s.posts = posts
	s.bySlug = bySlug
	s.mu.Unlock()
	return nil
}

// All returns all posts, ordered newest first
func (s *Store) All() []*Post {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Post, len(s.posts))
	copy(out, s.posts)
	return out
}

// Get returns the post with the given slug or nil
func (s *Store) Get(slug string) *Post {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.bySlug[slug]
}

// LatestUpdate returns the most recent post date in the store, or time.Time{} when empty. Used as Atom <updated>
func (s *Store) LatestUpdate() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var t time.Time
	for _, p := range s.posts {
		d := p.Meta.Updated
		if d.IsZero() {
			d = p.Meta.Date
		}
		if d.After(t) {
			t = d
		}
	}
	return t
}

// Splits a YAML front matter block (---) from the markdown body. Both parts might be empty
func parseFrontMatter(raw []byte) (FrontMatter, []byte, error) {
	var fm FrontMatter
	const delim = "---"
	s := string(raw)
	if !strings.HasPrefix(s, delim) {
		return fm, raw, nil
	}
	rest := s[len(delim):]
	rest = strings.TrimLeft(rest, "\r\n")
	end := strings.Index(rest, "\n"+delim)
	if end < 0 {
		return fm, raw, fmt.Errorf("unterminated front matter")
	}
	header := rest[:end]
	body := rest[end+len("\n"+delim):]
	body = strings.TrimLeft(body, "\r\n")

	if err := yaml.Unmarshal([]byte(header), &fm); err != nil {
		return fm, raw, err
	}
	return fm, []byte(body), nil
}
