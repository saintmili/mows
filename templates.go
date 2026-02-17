package mows

import (
	"html/template"
	"time"
)

// TemplateEngine wraps html/template
type TemplateEngine struct {
	pattern string
	funcMap template.FuncMap
	tmpl    *template.Template
}

// LoadTemplates loads templates using glob pattern.
// Example: "views/**/*.html"
func (e *Engine) LoadTemplates(pattern string) error {
	engine := &TemplateEngine{
		pattern: pattern,
		funcMap: defaultFuncMap(),
	}

	if err := engine.load(); err != nil {
		return err
	}

	e.templates = engine
	return nil
}

// defaultFuncMap returns the default template helper functions.
//
// These helpers are automatically available in all templates:
//
//   - safeHTML(string) → template.HTML : marks string as safe HTML
//   - now() → time.Time : current time
//   - date(time.Time, string) → string : formats time with layout
//
// Developers can also add custom functions via Engine.AddTemplateFunc.
func defaultFuncMap() template.FuncMap {
	return template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"now": func() time.Time {
			return time.Now()
		},
		"date": func(t time.Time, layout string) string {
			return t.Format(layout)
		},
	}
}

// AddTemplateFunc allows registering custom template helpers.
func (e *Engine) AddTemplateFunc(name string, fn any) {
	if e.templates == nil {
		e.templates = &TemplateEngine{
			funcMap: defaultFuncMap(),
		}
	}
	e.templates.funcMap[name] = fn
}

// load parses templates from the pattern and applies the template function map.
//
// If parsing fails, it returns an error. On success, the parsed templates
// are stored in the TemplateEngine.tmpl field.
//
// This function is called internally by Engine.LoadTemplates and automatically
// during DevMode hot-reload before rendering a template.
func (t *TemplateEngine) load() error {
	tmpl := template.New("").Funcs(t.funcMap)

	parsed, err := tmpl.ParseGlob(t.pattern)
	if err != nil {
		return err
	}

	t.tmpl = parsed
	return nil
}
