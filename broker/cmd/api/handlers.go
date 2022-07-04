package main

import (
	"broker/events"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
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
	case "log":
		app.publishLogEvent(w, requestPayload.Log)
	case "mail":
		app.sendEmail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) sendEmail(w http.ResponseWriter, mailData MailPayload) {
	jsonData, err := json.MarshalIndent(mailData, "", "\t")
	if err != nil {
		log.Println("Error trying to marshal mail data: ", err)
		app.errorJSON(w, err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, "http://mailer/send", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error trying create request to mailer: ", err)
		app.errorJSON(w, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		log.Println("Error on request to mailer: ", err)
		app.errorJSON(w, err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("Logger error response: ", err)
		app.errorJSON(w, err)
		return
	}

	responseJson := jsonResponse{
		Error:   false,
		Message: "Email sent successfully",
	}

	app.writeJSON(w, http.StatusOK, responseJson)
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	log.Println("Logging...")

	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		log.Println("Error trying to marshal log entry: ", err)
		app.errorJSON(w, err)
		return
	}

	request, err := http.NewRequest(http.MethodPost, "http://logger/log", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error trying create request to logger: ", err)
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Error on request to logger: ", err)
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		log.Println("Logger error response: ", err)
		app.errorJSON(w, err)
		return
	}

	responsePaylod := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(w, response.StatusCode, responsePaylod)
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

func (app *Config) publishLogEvent(w http.ResponseWriter, l LogPayload) {
	if err := app.pushToQueue(l.Name, l.Data); err != nil {
		log.Println(err)
		app.errorJSON(w, err)
	}

	responsePayload := jsonResponse{
		Error:   false,
		Message: "Logged via RabbitMQ",
	}

	app.writeJSON(w, http.StatusOK, responsePayload)
}

func (app *Config) pushToQueue(name, message string) error {
	emitter, err := events.NewEventEmitter(app.RabbitMQConn)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: message,
	}

	j, _ := json.Marshal(&payload)
	if err := emitter.Push(string(j), "log.INFO"); err != nil {
		return err
	}

	return nil
}
