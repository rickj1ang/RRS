package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthcheck", app.healthcheckHandler)
	mux.HandleFunc("GET /records", app.listRecordsHandler)
	mux.HandleFunc("POST /records", app.createRecordHandler)
	mux.HandleFunc("GET /records/{id}", app.showRecordHandler)
	mux.HandleFunc("DELETE /records/{id}", app.deleteRecordHandler)
	mux.HandleFunc("PATCH /records/{id}", app.updateRecordHandler)
	//user
	mux.HandleFunc("POST /users", app.registerUserHandler)
	mux.HandleFunc("GET /lord/{id}", app.givePowerHandler)

	//authentication
	mux.HandleFunc("POST /tokens/authentication", app.createAuthenticationTokenHandler)
	// this is a kind of custom Notfound page
	// if the URL do not match the routes up
	// stairs. app.notFoundResponse will out
	mux.HandleFunc("/", app.notFoundResponse)

	return app.recoverPanic(app.rateLimit(mux))
}
