package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
		return
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	log.Println("Trying to authenticate")
	jsonData, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request, err := http.NewRequest("POST", "http://auth/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}

	r, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer r.Body.Close()

	if r.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid crentials"), http.StatusUnauthorized)
		return
	} else if r.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	var reponseJson jsonResponse

	err = json.NewDecoder(r.Body).Decode(&reponseJson)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if reponseJson.Error {
		log.Printf("[ERROR]: reponseJson: %v", reponseJson)
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	var returnResponse jsonResponse
	returnResponse.Error = false
	returnResponse.Message = reponseJson.Message
	returnResponse.Data = reponseJson.Data

	app.writeJSON(w, http.StatusOK, returnResponse)
}
