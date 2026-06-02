// Provides a layer over html/template that mimics the ergonomics of Jinja2: a base layout with named blocks, partials that can be included from any page, and shared template functions
// Avoids larger libraries (pongo2, plush)
// html/template's context-aware escaping is the foundation; sanitization for blog HTML happens in internal/blog before reaching templates
package templates

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

// Renders named pages by composing a base layout with the page template and all available partials
type Engine struct {
	root  *template.Template
	pages map[string]*template.Template
	funcs template.FuncMap
}

// Configures the engine
type Options struct {
	// root templates directory containing layouts/, pages/, partials/
	Dir string
	// merged with the built-in function map
	Funcs template.FuncMap
	// base layout filename relative to Dir
	Layout string
}

// Parses all templates under opts.Dir into memory
// Returns an error if any template fails to parse
// Reload by calling New again, the engine itself is immutable after creation
func New(opts Options) (*Engine, error) {
	if opts.Dir == "" {
		opts.Dir = "web/templates"
	}
	if opts.Layout == "" {
		opts.Layout = "base.html"
	}

	funcs := defaultFuncs()
	for k, v := range opts.Funcs {
		funcs[k] = v
	}

	layoutPath := filepath.Join(opts.Dir, "layouts", opts.Layout)
	partials, err := filepath.Glob(filepath.Join(opts.Dir, "partials", "*.html"))
	if err != nil {
		return nil, err
	}
	pageFiles, err := filepath.Glob(filepath.Join(opts.Dir, "pages", "*.html"))
	if err != nil {
		return nil, err
	}

	pages := map[string]*template.Template{}
	for _, pf := range pageFiles {
		name := strings.TrimSuffix(filepath.Base(pf), ".html")
		t := template.New(filepath.Base(layoutPath)).Funcs(funcs)
		files := append([]string{layoutPath}, partials...)
		files = append(files, pf)
		parsed, err := t.ParseFiles(files...)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", name, err)
		}
		pages[name] = parsed
	}

	return &Engine{pages: pages, funcs: funcs}, nil
}

// Writes the named page using the base layout
// Data is exposed as the template's dot context
func (e *Engine) Render(w io.Writer, name string, data any) error {
	t, ok := e.pages[name]
	if !ok {
		return fmt.Errorf("template not found: %s", name)
	}
	return t.Execute(w, data)
}

// Reports whether a page template with the given name exists
func (e *Engine) Has(name string) bool { _, ok := e.pages[name]; return ok }

// Lists registered page templates
func (e *Engine) Names() []string {
	out := make([]string, 0, len(e.pages))
	for k := range e.pages {
		out = append(out, k)
	}
	return out
}

func defaultFuncs() template.FuncMap {
	return template.FuncMap{
		"year":     func() int { return time.Now().Year() },
		"fmtDate":  func(t time.Time) string { return t.Format("Jan 2, 2006") },
		"fmtISO":   func(t time.Time) string { return t.UTC().Format(time.RFC3339) },
		"upper":    strings.ToUpper,
		"lower":    strings.ToLower,
		"join":     strings.Join,
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
		"dict": func(values ...any) (map[string]any, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict requires even number of args")
			}
			m := make(map[string]any, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				k, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				m[k] = values[i+1]
			}
			return m, nil
		},
	}
}

// Helper that walks an fs.FS and returns all template file paths matching *.html. Good for embedded devices
func WalkTemplateFS(fsys fs.FS, root string) ([]string, error) {
	var out []string
	err := fs.WalkDir(fsys, root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(p, ".html") {
			out = append(out, p)
		}
		return nil
	})
	return out, err
}
