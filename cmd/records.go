package main

import (
	"fmt"
	"net/http"
)

func (app *application) createRecordHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create record")

}

func (app *application) showRecordHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDFromReq(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	//TBD: check id validation
	fmt.Fprintf(w, "show the record id = %d", id)

}
