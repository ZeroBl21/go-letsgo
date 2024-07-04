package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	base := CreateStack(app.recoverPanic, app.logRequest, secureHeader)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	return base(mux)
}
