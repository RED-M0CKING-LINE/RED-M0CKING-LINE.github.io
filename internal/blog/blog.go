// Loads, parses, and renders Markdown blog posts from the content directory
// Posts are loaded once on startup, so a restart is required for new content. This keeps it stateless
package blog

import (
	"bytes"
	"errors"
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
	alertcallouts "github.com/zmtcreative/gm-alert-callouts"
	"go.abhg.dev/goldmark/wikilink"
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

// Resolved WikiLinks
type resolvedWikilink struct{}

// Load scans dir for *.md files, parses front matter, renders HTML, and populates the store
// Drafts are skipped unless includeDrafts is true
func (s *Store) Load(dir string, includeDrafts bool) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read blog dir: %w", err)
	}

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

		clean, err := postRawToHTML(body)
		if err != nil {
			fmt.Errorf("render %s: %w", path, err)
		}

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

// Resolve WikiLinks to static blog content
func (r resolvedWikilink) ResolveWikilink(n *wikilink.Node) ([]byte, error) {
	target := string(n.Target)
	if n.Embed {
		// ![[attachments/foo.png]]
		return []byte("/static/assets/blog/" + target), nil
	}
	// [[Some Page]]
	if !strings.Contains(target, ".") {
		target += ".html"
	}
	return []byte("/static/assets/blog/" + target), nil
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

// Converts a raw markdown post to an HTML formatted post
func postRawToHTML(raw []byte) (template.HTML, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.Typographer,
			&wikilink.Extender{Resolver: resolvedWikilink{}},
			alertcallouts.NewAlertCallouts(
				alertcallouts.UseObsidianIcons(),
				alertcallouts.WithFolding(true),
				alertcallouts.WithCustomAlerts(true),
			),
		),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithHardWraps(), html.WithXHTML()),
	)
	// Strict sanitizing: allow Markdown, block scripts/iframes/etc
	policy := bluemonday.UGCPolicy()
	// Normal markdown formatting
	policy.AllowAttrs("class").OnElements("code", "pre", "span", "div")
	policy.AllowAttrs("id").OnElements("h1", "h2", "h3", "h4", "h5", "h6")
	// Wiki links and embedding images
	policy.AllowElements("img")
	policy.AllowAttrs("src", "alt").OnElements("img")
	// Callouts
	policy.AllowElements("details", "summary")
	policy.AllowAttrs("open").OnElements("details")
	policy.AllowAttrs("class", "data-callout").OnElements("details")
	policy.AllowAttrs("class", "data-callout").OnElements("div")
	// SVG icons emitted by gm-alert-callouts inside callout summary only
	policy.AllowElements("svg", "path", "g", "use", "defs", "symbol", "circle", "rect", "line", "polyline", "polygon")
	policy.AllowAttrs("xmlns", "viewBox", "fill", "stroke", "stroke-width", "stroke-linecap", "stroke-linejoin", "width", "height",
		"d", "cx", "cy", "r", "x", "y", "x1", "y1", "x2", "y2", "points", "transform", "aria-hidden").OnElements(
		"svg", "path", "g", "circle", "rect", "line", "polyline", "polygon")
	policy.AllowAttrs("class").OnElements(
		"svg", "path", "g", "circle", "rect", "line", "polyline", "polygon", "summary")

	var buf bytes.Buffer
	if err := md.Convert(raw, &buf); err != nil {
		return "", errors.New("error rendering HTML")
	}
	clean := policy.SanitizeBytes(buf.Bytes())
	return template.HTML(clean), nil
}

// Truncate the post content string based on line OR character count, and append a string when truncation has occured
func truncatePostContent(p Post, truncatedAppend string, limitChars uint, limitLines uint) (Post, error) {
	if (limitChars == 0 && limitLines == 0) || (limitChars != 0 && limitLines != 0) {
		return p, errors.New("limitChars xor limitLines must be non zero")
	}

	if limitChars != 0 {
		runes := []rune(p.Raw)
		if uint(len(runes)) > limitChars {
			p.Raw = string(runes[:limitChars]) + truncatedAppend
		}
	}

	if limitLines != 0 {
		lines := strings.Split(p.Raw, "\n")
		if uint(len(lines)) > limitLines {
			p.Raw = strings.Join(lines[:limitLines], "\n") + truncatedAppend
		}
	}
	clean, err := postRawToHTML([]byte(p.Raw))
	if err != nil {
		errors.New("error converting new raw to html")
	}
	p.HTML = clean

	return p, nil
}
