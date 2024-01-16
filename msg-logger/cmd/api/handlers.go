package main

import (
	"msg-logger/data"
	"net/http"
)

type JSONPayLoad struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayLoad JSONPayLoad
	_ = app.readJSON(w, r, &requestPayLoad)

	event := data.LogEntry{
		Name: requestPayLoad.Name,
		Data: requestPayLoad.Data,
	}

	err := app.Models.LogEntry.Insert(event)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	response := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.WriteJSON(w, http.StatusAccepted, response)
}
