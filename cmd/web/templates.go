package main

import (
	"github.com/ganeshrahul23/snippetbox/pkg/forms"
	"github.com/ganeshrahul23/snippetbox/pkg/models"
	"html/template"
	"path/filepath"
	"time"
)

type templateData struct {
	AuthenticatedUser *models.User
	CurrentYear       int
	Snippet           *models.Snippet
	Snippets          []*models.Snippet
	Form              *forms.Form
	Flash             string
	CSRFToken         string
}

func humanDate(t time.Time) string {
	//return t.Format("02 Jan 2006 at 15:04")

	// Return the empty string if time has the zero value.
	if t.IsZero() {
		return ""
	}
	// Convert the time to UTC before formatting it.
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))

	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		//ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
