package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

	err = app.logRequest("Auth app", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		log.Printf("error trying to log authentication:  %v", err)
		app.errorJSON(w, err, http.StatusBadRequest)
	}

	responsePayload := jsonResponse{
		Error:   false,
		Message: "Authentication successfull",
		Data:    user,
	}

	app.writeJSON(w, http.StatusOK, responsePayload)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	j, _ := json.MarshalIndent(entry, "", "\t")
	request, err := http.NewRequest(http.MethodPost, "http://logger/log", bytes.NewBuffer(j))
	if err != nil {
		return err
	}

	client := &http.Client{}

	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
