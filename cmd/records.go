package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rickj1ang/RRS/internal/data"
	"github.com/rickj1ang/RRS/internal/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		record.LastChange = time.Now()
	}
	if input.Description != nil {
		record.Description = *input.Description
	}
	if input.Genres != nil {
		record.Genres = input.Genres
	}
	v := validator.New()

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

func (app *application) listAllRecordsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	res, err := app.models.Records.GetAll(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusFound, envelope{"records": res}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) readBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDFromReq(r)
	if err != nil || id == primitive.NilObjectID {
		app.badRequestResponse(w, r, err)
		return
	}

	page, err := readPageFromReq(r)
	if err != nil || page == -1 {
		app.badRequestResponse(w, r, err)
		return
	}

	record, err := app.models.Records.Get("_id", id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(record.TotalPages >= uint16(page), "page", "current page can not bigger than total pages")
	v.Check(page >= 0, "page", "page must bigger than 0")
	if !v.Valid() {
		app.failValidationResponse(w, r, v.Errors)
		return
	}

	record.CurrentPage = uint16(page)
	record.Progress = float32(record.CurrentPage) / float32(record.TotalPages)

	err = app.models.Records.Update(id, record)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusCreated, envelope{"update": fmt.Sprintf("Now you read to %f of the %s", record.Progress, record.Title)}, nil)
}
