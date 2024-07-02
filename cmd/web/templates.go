package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ZeroBl21/go-letsgo/internal/models"
)

// Holds the structure for any dynamicfor HTML template
type templateData struct {
	CurrentYear int

	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func (app *application) newTemplateData(_ *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
