package main

import (
	"log"
	"logger/data"
	"net/http"
)

type RequestPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	_ = app.readJSON(w, r, &payload)

	event := data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		log.Panicln("WriteLog ERROR: ", err)
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Message: "Logged",
		Error:   false,
	}

	app.writeJSON(w, http.StatusCreated, resp)
}
