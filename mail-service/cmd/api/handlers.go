package main

import (
	"log"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {

	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var request mailMessage

	err := app.readJSON(w, r, &request)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	msg := Message{
		From:    request.From,
		To:      request.To,
		Subject: request.Subject,
		Data:    request.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	log.Printf("Error sending mail is %s", err)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: "Message is sent to" + request.To,
	}

	app.writeJSON(w, http.StatusOK, payload)

}
