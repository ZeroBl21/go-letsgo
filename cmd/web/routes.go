package main

import (
	"net/http"

	"github.com/ZeroBl21/go-letsgo/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	base := CreateStack(app.recoverPanic, app.logRequest, secureHeader)
	dynamic := CreateStack(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	protected := CreateStack(dynamic, app.requireAuthentication)

	fileServer := http.FileServer(http.FS(ui.Files))
	mux.Handle("GET /static/", fileServer)

	// Home
	mux.HandleFunc("/{$}", dynamic.ToHandlerFunc(app.home))

	// Snippets
	mux.HandleFunc("GET /snippet/view/{id}", dynamic.ToHandlerFunc(app.snippetView))

	// Protected
	mux.HandleFunc("GET /snippet/create", protected.ToHandlerFunc(app.snippetCreate))
	mux.HandleFunc("POST /snippet/create", protected.ToHandlerFunc(app.snippetCreatePost))

	// Auth
	mux.HandleFunc("GET /user/register", dynamic.ToHandlerFunc(app.userRegister))
	mux.HandleFunc("POST /user/register", dynamic.ToHandlerFunc(app.userRegisterPost))

	mux.HandleFunc("GET /user/login", dynamic.ToHandlerFunc(app.userLogin))
	mux.HandleFunc("POST /user/login", dynamic.ToHandlerFunc(app.userLoginPost))

	// Protected
	mux.HandleFunc("POST /user/logout", protected.ToHandlerFunc(app.userLogoutPost))

	return base(mux)
}
