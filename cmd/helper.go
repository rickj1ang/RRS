package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func readIDFromReq(r *http.Request) (int64, error) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	if id < 1 {
		return 0, errors.New("id Can not less than 1")
	}
	return id, nil
}

// json write helper
// data as body set status headers as header write into w
func writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	//TBD: change a package to handle JSON for higher perfermance
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, val := range headers {
		w.Header()[key] = val
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
