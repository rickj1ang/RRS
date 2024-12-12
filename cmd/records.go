package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rickj1ang/RRS/internal/data"
)

func (app *application) createRecordHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create record")

}

func (app *application) showRecordHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDFromReq(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	record := data.Record{
		ID:          id,
		CreatedAt:   time.Now(),
		Title:       "Republic",
		Writer:      "Plato",
		TotalPages:  400,
		CurrentPage: 200,
		Description: "GOod stuff",
		Genres:      []string{"philo", "conversation"},
	}

	record.Progress = float32(record.CurrentPage) / float32(record.TotalPages)

	//TBD: check id validation
	err = writeJSON(w, http.StatusOK, envelope{"record": record}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
