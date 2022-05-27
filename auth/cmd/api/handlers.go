package main

import (
	"errors"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Printf("[ERROR]: %v", err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		log.Println("User found: ", user)
		app.writeJSON(w, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		log.Printf("[ERROR]: %v\n\tVALID: %v", err, valid)
		app.writeJSON(w, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	responsePayload := jsonResponse{
		Error:   false,
		Message: "Authentication successfull",
		Data:    user,
	}

	app.writeJSON(w, http.StatusOK, responsePayload)
}
