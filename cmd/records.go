package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rickj1ang/RRS/internal/data"
)

func (app *application) createRecordHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string   `json:"title"`
		Writer      string   `json:"writer,omitempty"`
		TotalPages  uint16   `json:"total_pages,omitempty"`
		CurrentPage uint16   `json:"curent_page,omitempty"`
		Description string   `json:"description,omitempty"`
		Genres      []string `json:"genres,omitempty"`
	}

	err := readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
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
