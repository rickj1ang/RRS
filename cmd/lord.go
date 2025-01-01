package main

import (
	"net/http"
)

func (app *application) givePowerHandler(w http.ResponseWriter, r *http.Request) {
	//TBD: Lord check! only lord can use this function

	id, err := readIDFromReq(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.Get("_id", id)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user.Level = 1

	err = app.models.Users.Update("_id", id, user)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// TBD: check id validation
	err = writeJSON(w, http.StatusOK, envelope{"userAfterChange": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
