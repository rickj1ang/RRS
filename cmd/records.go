package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rickj1ang/RRS/internal/data"
	"github.com/rickj1ang/RRS/internal/validator"
)

func (app *application) createRecordHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string   `json:"title"`
		Writer      string   `json:"writer,omitempty"`
		TotalPages  uint16   `json:"total_pages,omitempty"`
		CurrentPage uint16   `json:"current_page,omitempty"`
		Description string   `json:"description,omitempty"`
		Genres      []string `json:"genres,omitempty"`
	}

	err := readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	record := &data.Record{
		Title:       input.Title,
		Writer:      input.Writer,
		TotalPages:  input.TotalPages,
		CurrentPage: input.CurrentPage,
		Description: input.Description,
		Genres:      input.Genres,
	}
	record.Progress = float32(record.CurrentPage) / float32(record.TotalPages)
	record.CreatedAt = time.Now()

	if data.ValidateRecord(v, record); !v.Valid() {
		app.failValidationResponse(w, r, v.Errors)
		return
	}

	insertId, err := app.models.Records.Insert(record)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	app.logger.info.Printf("Insert a piece of document which id is %s", insertId)

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("records/%s", insertId))

	err = writeJSON(w, http.StatusCreated, envelope{"record": record}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showRecordHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDFromReq(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	record, err := app.models.Records.Get("_id", id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//TBD: check id validation
	err = writeJSON(w, http.StatusOK, envelope{"record": record}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
