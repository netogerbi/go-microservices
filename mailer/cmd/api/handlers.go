package main

import (
	"fmt"
	"log"
	"net/http"
)

func (app *Config) sendMail(w http.ResponseWriter, r *http.Request) {
	var mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	if err := app.readJSON(w, r, &mailMessage); err != nil {
		log.Println("Error trying to read JSON: ", err)
		app.errorJSON(w, err)
	}

	msg := Message{
		From:    mailMessage.From,
		To:      mailMessage.To,
		Subject: mailMessage.Subject,
		Data:    mailMessage.Message,
	}

	if err := app.Mailer.SendSMTPMessage(msg); err != nil {
		log.Println("Error trying to send email: ", err)
		app.errorJSON(w, err)
	}

	res := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Mail sent to %s", msg.To),
	}

	app.writeJSON(w, http.StatusOK, res)
}
