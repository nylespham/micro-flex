package main

import "net/http"

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}
	var requestPayLoad mailMessage

	err := app.readJSON(w, r, &requestPayLoad)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	msg := Message{
		From:    requestPayLoad.From,
		To:      requestPayLoad.To,
		Subject: requestPayLoad.Subject,
		Data:    requestPayLoad.Message,
	}

	err = app.Mailer.sendSMTPMessage(msg)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payLoad := jsonResponse{
		Error:   false,
		Message: "sent to " + requestPayLoad.To,
	}

	app.WriteJSON(w, http.StatusAccepted, payLoad)
}
