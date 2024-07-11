package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	base := CreateStack(app.recoverPanic, app.logRequest, secureHeader)
	dynamic := CreateStack(app.sessionManager.LoadAndSave)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/{$}", dynamic.ToHandlerFunc(app.home))
	mux.HandleFunc("GET /snippet/view/{id}", dynamic.ToHandlerFunc(app.snippetView))
	mux.HandleFunc("GET /snippet/create", dynamic.ToHandlerFunc(app.snippetCreate))
	mux.HandleFunc("POST /snippet/create", dynamic.ToHandlerFunc(app.snippetCreatePost))

	return base(mux)
}
