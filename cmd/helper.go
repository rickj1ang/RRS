package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

type envelope map[string]any

// json write helper
// data as body set status headers as header write into w
func writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	//TBD: change a package to handle JSON for higher perfermance
	// MarshalIndent has lower perfermance than Marshal.
	js, err := json.MarshalIndent(data, "", "\t")
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

func readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// JSON syntax is wrong
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		// same as io.EOF but in the middle
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("badly contains badly-formed JSON")

		// the values not suit the GO struct we want to put them
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// request body is io.EOF(nil)
		case errors.Is(err, io.EOF):
			return errors.New("Body must have something")

		// we pass a no-nil pointer to Decode()
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		// other errors we can not predict
		default:
			return err
		}
	}
	return nil
}
