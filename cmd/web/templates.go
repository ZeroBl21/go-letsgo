package main

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ZeroBl21/go-letsgo/internal/models"
	"github.com/ZeroBl21/go-letsgo/ui"
	"github.com/justinas/nosurf"
)

// Holds the structure for any dynamic for HTML template
type templateData struct {
	CurrentYear int

	Snippet  *models.Snippet
	Snippets []*models.Snippet

	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		ts, err := template.
			New(name).
			Funcs(functions).
			ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

// "V-Table" for custom templates functions
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// Returns a nicely formatted string representation of time.Time object
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 3:04pm")
}
