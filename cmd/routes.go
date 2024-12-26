package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthcheck", app.healthcheckHandler)
	mux.HandleFunc("GET /records", app.listRecordsHandler)
	mux.HandleFunc("POST /records", app.createRecordHandler)
	mux.HandleFunc("GET /records/{id}", app.showRecordHandler)
	mux.HandleFunc("DELETE /records/{id}", app.deleteRecordHandler)
	mux.HandleFunc("PATCH /records/{id}", app.updateRecordHandler)
	// this is a kind of custom Notfound page
	// if the URL do not match the routes up
	// stairs. app.notFoundResponse will out
	mux.HandleFunc("/", app.notFoundResponse)

	return mux
}
