package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Writesa an error message and the stack trace to the errorLog, then sends
// a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

// Sends a specific status code and corresponding description to the user.
func (app *application) clienError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Send 404 not found to the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clienError(w, http.StatusNotFound)
}
