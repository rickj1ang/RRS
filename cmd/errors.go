package main

import (
	"fmt"
	"net/http"
)

// print to terminal for developer to debug, can write to a logfile
func (app *application) logError(r *http.Request, err error) {
	app.loggerError.Println(err)
}

// response to API client
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request,
	status int, message any) {
	env := envelope{"error": message}
	err := writeJSON(w, status, env, nil)
	// It's a err of err wrtiing, a little bit funny :)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// this error we need to know so must log it, but notfound and method not allow
	// is what we plan it, so just client need to know, do not need to log it
	app.logError(r, err)

	message := "Server encountered a problem, can not process your Request\n"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "The resource you are requesting can not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not support for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
