package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayLoad struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &requestPayLoad)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	// validate user against the database
	user, err := app.Models.User.GetByEmail(requestPayLoad.Email)

	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayLoad.Password)

	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}
	// log authentication

	err = app.logRequest("authenticate", fmt.Sprintf("user %s authenticated", user.Email))

	payLoad := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.WriteJSON(w, http.StatusAccepted, payLoad)

}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceUrl := "http://msg-logger/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))

	if err != nil {
		log.Println(err)
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)

	if err != nil {
		return err
	}

	return nil
}
