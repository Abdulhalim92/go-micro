package main

import (
	"log"
	"net/http"
)

type MailMessagePayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// SendMail - обработчик для отправки почты через SMTP.
func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	var requestPayload MailMessagePayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Printf("Error reading request body: %s", err)
		app.errorJSON(w, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SentSMTPMessage(msg)
	if err != nil {
		log.Println("Error sending mail:", err)
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
