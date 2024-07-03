package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	base := CreateStack(app.logRequest, secureHeader)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/{$}", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("POST /snippet/create", app.snippetCreate)

	return base(mux)
}
