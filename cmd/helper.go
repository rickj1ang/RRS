package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/rickj1ang/RRS/internal/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func readIDFromReq(r *http.Request) (primitive.ObjectID, error) {
	idStr := r.PathValue("id")

	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return objectID, nil

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
	// we set a limit body size 1M it's enough for every record
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	// if anything unexpected appear in JSON, it will return error
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
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

		// detect error if it has any wnknown key in JSON it will point it out
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// the request body larger than 1MB
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must be less than %d bytes", maxBytes)

		// other errors we can not predict
		default:
			return err
		}

	}
	// we already Decode a JSON value, if we can abstract something
	// now, it means Bad Request which contains more than one JSON value
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contains a single JSON value")

	}
	return nil
}

func (app *application) readString(qs url.Values, key, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer vale")
		return defaultValue
	}
	return i
}

// run fn bakcground with panic recover
func (app *application) background(fn func()) {
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()

		defer func() {

			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()

		fn()
	}()
}
