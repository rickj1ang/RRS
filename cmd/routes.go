package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthcheck", app.healthcheckHandler)
	mux.HandleFunc("POST /records", app.createRecordHandler)
	mux.HandleFunc("GET /records/{id}", app.showRecordHandler)

	return mux
}
