package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"msg-broker/event"
	"net/http"
	"net/rpc"
)

type RequestPayLoad struct {
	Action string      `json:"action"`
	Auth   AuthPayLoad `json:"auth,omitempty"`
	Log    LogPayLoad  `json:"log,omitempty"`
	Mail   MailPayLoad `json:"mail,omitempty"`
}

type MailPayLoad struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AuthPayLoad struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayLoad struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payLoad := jsonResponse{
		Error:   false,
		Message: "hit the broker",
	}
	_ = app.WriteJSON(w, http.StatusOK, payLoad)
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayLoad) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceUrl := "http://msg-logger:12800/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse

	payload.Error = false
	payload.Message = "Successfully logged item"

	app.WriteJSON(w, http.StatusAccepted, payload)
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
	case "log":
		app.logItemViaRPC(w, requestPayLoad.Log)
	case "mail":
		app.sendMail(w, requestPayLoad.Mail)
	case "rabbit":
		app.logEventViaRabbit(w, requestPayLoad.Log)
	default:
		app.errorJSON(w, errors.New("invalid action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayLoad) {
	// create some JSON that we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the login-oauth microservice
	request, err := http.NewRequest("POST", "http://login-oauth:13780/authenticate", bytes.NewBuffer(jsonData))
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
	} else if response.StatusCode != http.StatusAccepted {
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

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayLoad) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	mailServiceUrl := "http://msg-mail:4100/send"

	request, err := http.NewRequest("POST", mailServiceUrl, bytes.NewBuffer(jsonData))

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mail microservice"))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Successfully sent mail"

	app.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayLoad) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Successfully logged item"
	app.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(&app.Rabbit)
	if err != nil {
		return err
	}
	payLoad := LogPayLoad{
		Name: name,
		Data: msg,
	}
	j, _ := json.MarshalIndent(&payLoad, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}

type RPCPayLoad struct {
	Name string
	Data string
}

func (app *Config) logItemViaRPC(w http.ResponseWriter, l LogPayLoad) {
	client, err := rpc.Dial("tcp", "msg-logger:12800")

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	rpcPayLoad := RPCPayLoad{
		Name: l.Name,
		Data: l.Data,
	}

	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayLoad, &result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	payLoad := jsonResponse{
		Error:   false,
		Message: result,
	}

	app.WriteJSON(w, http.StatusAccepted, payLoad)

}
