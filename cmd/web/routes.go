package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	base := CreateStack(app.recoverPanic, app.logRequest, secureHeader)
	dynamic := CreateStack(app.sessionManager.LoadAndSave)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Home
	mux.HandleFunc("/{$}", dynamic.ToHandlerFunc(app.home))

	// Snippets
	mux.HandleFunc("GET /snippet/view/{id}", dynamic.ToHandlerFunc(app.snippetView))
	mux.HandleFunc("GET /snippet/create", dynamic.ToHandlerFunc(app.snippetCreate))
	mux.HandleFunc("POST /snippet/create", dynamic.ToHandlerFunc(app.snippetCreatePost))

	// Auth
	mux.HandleFunc("GET /user/register", dynamic.ToHandlerFunc(app.userRegister))
	mux.HandleFunc("POST /user/register", dynamic.ToHandlerFunc(app.userRegisterPost))

	mux.HandleFunc("GET /user/login", dynamic.ToHandlerFunc(app.userLogin))
	mux.HandleFunc("POST /user/login", dynamic.ToHandlerFunc(app.userLoginPost))

	mux.HandleFunc("POST /user/logout", dynamic.ToHandlerFunc(app.userLogoutPost))

	return base(mux)
}
