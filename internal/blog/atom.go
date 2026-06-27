package blog

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Atom feed types per RFC 4287 https://www.rfc-editor.org/info/rfc4287/
type atomFeed struct {
	XMLName  xml.Name    `xml:"http://www.w3.org/2005/Atom feed"`
	Title    string      `xml:"title"`
	Subtitle string      `xml:"subtitle,omitempty"`
	ID       string      `xml:"id"`
	Updated  string      `xml:"updated"`
	Links    []atomLink  `xml:"link"`
	Author   *atomAuthor `xml:"author,omitempty"`
	Entries  []atomEntry `xml:"entry"`
}

type atomLink struct {
	Rel  string `xml:"rel,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
	Href string `xml:"href,attr"`
}

type atomAuthor struct {
	Name string `xml:"name"`
}

type atomEntry struct {
	Title      string         `xml:"title"`
	ID         string         `xml:"id"`
	Link       atomLink       `xml:"link"`
	Published  string         `xml:"published,omitempty"`
	Updated    string         `xml:"updated"`
	Summary    string         `xml:"summary,omitempty"`
	Content    atomContent    `xml:"content"`
	Author     *atomAuthor    `xml:"author,omitempty"`
	Categories []atomCategory `xml:"category,omitempty"`
}

type atomContent struct {
	Type string `xml:"type,attr"`
	Body string `xml:",chardata"`
}

type atomCategory struct {
	Term string `xml:"term,attr"`
}

// Defaults
type FeedOptions struct {
	SiteName string
	BaseURL  string
	Author   string
}

// AtomFeed serializes the store to an Atom 1.0 XML document
func (s *Store) AtomFeed(opt FeedOptions) ([]byte, error) {
	base := strings.TrimRight(opt.BaseURL, "/")
	updated := s.LatestUpdate()
	if updated.IsZero() {
		updated = time.Now().UTC()
	}

	feed := atomFeed{
		Title:    opt.SiteName,
		Subtitle: "Notes From Infrastructure Trenches",
		ID:       base + "/",
		Updated:  updated.UTC().Format(time.RFC3339),
		Links: []atomLink{
			{Rel: "self", Type: "application/atom+xml", Href: base + "/feed.xml"},
			{Rel: "alternate", Type: "text/html", Href: base + "/blog"},
		},
	}
	if opt.Author != "" {
		feed.Author = &atomAuthor{Name: opt.Author}
	}

	for _, p := range s.All() {
		entryURL := fmt.Sprintf("%s/blog/%s", base, p.Slug)
		published := p.Meta.Date
		upd := p.Meta.Updated
		if upd.IsZero() {
			upd = published
		}
		if upd.IsZero() {
			upd = p.ModTime
		}

		e := atomEntry{
			Title:     p.Meta.Title,
			ID:        entryURL,
			Link:      atomLink{Rel: "alternate", Type: "text/html", Href: entryURL},
			Published: rfc3339OrEmpty(published),
			Updated:   rfc3339OrEmpty(upd),
			Summary:   p.Meta.Summary,
		}

		if p.Meta.Author != "" {
			e.Author = &atomAuthor{Name: p.Meta.Author}
		}
		for _, t := range p.Meta.Tags {
			e.Categories = append(e.Categories, atomCategory{Term: t})
		}
		truncatedPostContent, err := truncatePostContent(p, "...\nContent truncated. Continue reading on the website.", 3000, 0)
		if err != nil {
			return nil, err
		}
		e.Content = atomContent{Type: "html", Body: string(truncatedPostContent.HTML)}
		feed.Entries = append(feed.Entries, e)
	}

	var buf strings.Builder
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(&strBuilderWriter{&buf})
	enc.Indent("", "  ")
	if err := enc.Encode(feed); err != nil {
		return nil, err
	}
	if err := enc.Flush(); err != nil {
		return nil, err
	}
	return []byte(buf.String()), nil
}

func rfc3339OrEmpty(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

// adapts strings.Builder to io.Writer
type strBuilderWriter struct{ b *strings.Builder }

func (w *strBuilderWriter) Write(p []byte) (int, error) { return w.b.Write(p) }

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

// Truncate the post content string based on line OR character count, and append a string when truncation has occured
func truncatePostContent(p *Post, truncatedAppend string, limitChars uint, limitLines uint) (*Post, error) {
	if (limitChars == 0 && limitLines == 0) || (limitChars != 0 && limitLines != 0) {
		return nil, errors.New("limitChars xor limitLines must be non zero")
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
