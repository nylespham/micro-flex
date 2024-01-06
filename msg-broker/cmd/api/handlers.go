package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payLoad := jsonResponse{
		Error:   false,
		Message: "hit the broker",
	}
	_ = app.WriteJSON(w, http.StatusAccepted, payLoad)
}
