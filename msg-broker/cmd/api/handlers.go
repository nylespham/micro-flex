package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayLoad struct {
	Action string      `json:"action"`
	Auth   AuthPayLoad `json:"auth,omitempty"`
}

type AuthPayLoad struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payLoad := jsonResponse{
		Error:   false,
		Message: "hit the broker",
	}
	_ = app.WriteJSON(w, http.StatusAccepted, payLoad)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayLoad RequestPayLoad
	err := app.readJSON(w, r, &requestPayLoad)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayLoad.Action {
	case "auth":
		app.authenticate(w, requestPayLoad.Auth)

	default:
		app.errorJSON(w, errors.New("invalid action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayLoad) {
	// create some JSON that we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the login-oauth microservice
	request, err := http.NewRequest("POST", "http://localhost:1800/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode == http.StatusAccepted {
		app.errorJSON(w, errors.New("errors calling auth microservice"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	// decode json from auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payLoad jsonResponse
	payLoad.Error = false
	payLoad.Message = "Succesfully authenticated"
	payLoad.Data = jsonFromService.Data

	app.WriteJSON(w, http.StatusAccepted, payLoad)
}
