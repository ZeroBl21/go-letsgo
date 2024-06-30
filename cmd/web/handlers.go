package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, _ *http.Request) {
	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/home.html",

		// Partials
		"./ui/html/partials/nav.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err := ts.ExecuteTemplate(w, "base", nil); err != nil {
		app.serverError(w, err)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Create a new snippet...")
}
