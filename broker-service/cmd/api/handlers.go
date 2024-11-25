package main

import (
	"broker/event"
	"broker/logs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	FromAddress string `json:"from"`
	ToAddress   string `json:"to"`
	Subject     string `json:"subject"`
	Message     string `json:"message"`
}
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse{
		Message: "Hit the broker",
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	var request RequestPayload
	err := app.readJSON(w, r, &request)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	conn, err := grpc.NewClient("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer conn.Close()

	c := logs.NewLogServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	logPayload := request.Log
	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: logPayload.Name,
			Data: logPayload.Data,
		},
	})

	if err != nil {
		app.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Message: "Logged",
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var request RequestPayload
	err := app.readJSON(w, r, &request)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	switch request.Action {
	case "auth":
		app.authenticate(w, request.Auth)
	case "log":
		app.logItemViaRPC(w, request.Log)
	case "mail":
		app.SendMail(w, request.Mail)
	default:
		app.errorJson(w, errors.New("unknown action"))
	}

}

func (app *Config) authenticate(w http.ResponseWriter, payload AuthPayload) {
	jsonData, _ := json.MarshalIndent(payload, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusOK {
		app.errorJson(w, errors.New("error calling auth service"), response.StatusCode)
		return
	}

	var res jsonResponse

	err = json.NewDecoder(request.Body).Decode(&res)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	if res.Error {
		app.errorJson(w, err)
		return
	}

	var resp jsonResponse
	resp.Error = true
	resp.Message = "Authenticated"
	resp.Data = res
	app.writeJSON(w, http.StatusOK, resp)

}

func (app *Config) logItem(w http.ResponseWriter, payload LogPayload) {

	jsonData, _ := json.Marshal(payload)

	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		app.errorJson(w, errors.New("error calling log service"), response.StatusCode)
		return
	}

	var resp jsonResponse
	resp.Error = true
	resp.Message = "Logged"
	app.writeJSON(w, http.StatusOK, resp)

}
func (app *Config) logItemViaRabbitmq(w http.ResponseWriter, payload LogPayload) {
	err := app.pushToQueue(payload.Name, payload.Data)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	var res jsonResponse
	res.Error = false
	res.Message = "Logged via Rabitmq"

	app.writeJSON(w, http.StatusOK, res)

}

func (app *Config) pushToQueue(name, message string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return nil
	}

	payload := LogPayload{
		Name: name,
		Data: message,
	}
	j, _ := json.Marshal(&payload)
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return nil
	}
	return nil
}

type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) logItemViaRPC(w http.ResponseWriter, payload LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")

	if err != nil {
		app.errorJson(w, err)
		return
	}

	rpcPayload := RPCPayload{
		Name: payload.Name,
		Data: payload.Data,
	}

	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	res := jsonResponse{
		Error:   false,
		Message: result,
	}

	app.writeJSON(w, http.StatusOK, res)

}

func (app *Config) SendMail(w http.ResponseWriter, payload MailPayload) {

	jsonData, _ := json.Marshal(payload)

	request, err := http.NewRequest("POST", "http://mail-service/send", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.errorJson(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		app.errorJson(w, errors.New("error calling mail service"), response.StatusCode)
		return
	}
	var resp jsonResponse
	resp.Error = true
	resp.Message = "Mail sent"
	app.writeJSON(w, http.StatusOK, resp)
}
