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

	user := app.contextGetUser(r)
	record := &data.Record{
		Title:       input.Title,
		Writer:      input.Writer,
		TotalPages:  input.TotalPages,
		CurrentPage: input.CurrentPage,
		Description: input.Description,
		Genres:      input.Genres,
		Owner:       user.ID,
	}
	record.Progress = float32(record.CurrentPage) / float32(record.TotalPages)
	record.LastChange = time.Now()

	if data.ValidateRecord(v, record); !v.Valid() {
		app.failValidationResponse(w, r, v.Errors)
		return
	}

	insertId, err := app.models.Records.Insert(record)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	user.Records = append(user.Records, insertId)
	err = app.models.Users.Update("_id", user.ID, user)

	app.logger.PrintInfo(fmt.Sprintf("Insert a piece of document which id is %s", insertId.Hex()), nil)

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("records/%s", insertId.Hex()))

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

func (app *application) deleteRecordHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDFromReq(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	n, err := app.models.Records.Delete(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	message := fmt.Sprintf("successful delete %d item", n)

	err = writeJSON(w, http.StatusOK, envelope{"message": message}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// I do not concern race-condition in this function,
// It make no sense for a user to login in two device
// and make two request concurrency but I will do it later
func (app *application) updateRecordHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDFromReq(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input struct {
		Title       *string  `json:"title"`
		Writer      *string  `json:"writer"`
		TotalPages  *uint16  `json:"total_pages"`
		CurrentPage *uint16  `json:"current_page"`
		Description *string  `json:"description"`
		Genres      []string `json:"genres"`
	}

	err = readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	record, err := app.models.Records.Get("_id", id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if input.Title != nil {
		record.Title = *input.Title
	}
	if input.Writer != nil {
		record.Writer = *input.Writer
	}
	if input.TotalPages != nil {
		record.TotalPages = *input.TotalPages
	}
	if input.CurrentPage != nil {
		record.CurrentPage = *input.CurrentPage
		record.Progress = float32(record.CurrentPage) / float32(record.TotalPages)
	}
	if input.Description != nil {
		record.Description = *input.Description
	}
	if input.Genres != nil {
		record.Genres = input.Genres
	}
	v := validator.New()

	record.LastChange = time.Now()
	if data.ValidateRecord(v, record); !v.Valid() {
		app.failValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Records.Update(id, record)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = writeJSON(w, http.StatusOK, envelope{"record": record}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listRecordsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 5, v)

	input.Filters.Sort = app.readString(qs, "sort", "created_at")
	input.Filters.SortSafelist = []string{"created_at", "title", "year", "runtime", "-created_at", "-title", "-year", "-runtime"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}
