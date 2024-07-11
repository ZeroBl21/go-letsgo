package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ZeroBl21/go-letsgo/internal/models"
	"github.com/ZeroBl21/go-letsgo/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.html", data)
}

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clienError(w, http.StatusBadRequest)
		return
	}

	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clienError(w, http.StatusBadRequest)
		return
	}

	// Title
	form.CheckField(validator.NotBlank(form.Title), "title",
		"This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title",
		"This field cannot be be more than 100 characters long")

	// Content
	form.CheckField(validator.NotBlank(form.Content), "content",
		"This field cannot be blank")

	// Expires
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires",
		"This field must be equal to 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
