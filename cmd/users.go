package main

import (
	"errors"
	"net/http"

	"github.com/rickj1ang/RRS/internal/data"
	"github.com/rickj1ang/RRS/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	err := readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:  input.Name,
		Email: input.Email,
	}

	user.Level = 0
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failValidationResponse(w, r, v.Errors)
		return
	}

	userID, err := app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "This email is already used")
			app.failValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	message := "send your pure userID to rick0j1ang@gmail.com with the email address you used to sign up"
	err = writeJSON(w, http.StatusCreated, envelope{"userID": userID, "activation": message}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
