package main

import (
	"broker/event"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// RequestPayload - структура для передачи данных от клиента
type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

// AuthPayload - структура для передачи данных аутентификации
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// MailPayload - структура для передачи данных почты
type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// LogPayload - структура для передачи логов
type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// HandleSubmission - обрабатывает запросы от клиента и перенаправляет их в нужный сервис
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logEventViaRabbit(w, requestPayload.Log)
		//app.logItem(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

// logEventViaRabbit - отправляет событие в RabbitMQ
func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = "logged via RabbitMQ"

	app.writeJSON(w, http.StatusAccepted, payload)
}

// pushToQueue - отправляет сообщение в очередь RabbitMQ
func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.RabbitMQ)
	if err != nil {
		log.Printf("Error creating event emitter: %s", err)
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	jsonData, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		log.Printf("Error marshalling payload: %s", err)
		return err
	}

	err = emitter.Push(string(jsonData), "log.INFO")
	if err != nil {
		log.Printf("Error pushing to queue: %s", err)
		return err
	}

	return nil
}

// sendMail - отправляет запрос на отправку почты в сервис почты
func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// call the mail service
	mailServiceURL := "http://mailer-service/send"

	// post to mail service
	request, err := http.NewRequest(http.MethodPost, mailServiceURL, bytes.NewBuffer(jsonData))
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

	// make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mail service"))
		return
	}

	// send back json
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	app.writeJSON(w, http.StatusAccepted, payload)
}

// logItem - отправляет лог в сервис логирования
func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating log request: %v", err)
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error creating log request: %v", err)
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("bad status code"))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(w, http.StatusAccepted, payload)
}

// authenticate - отправляет запрос на аутентификацию в сервис аутентификации
func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest(http.MethodPost, "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[authenticate] - > error creating request: %v", err)
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Printf("[autneticate] -> error calling authentication service: %v", err)
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	log.Printf("[authenticate] - > authentication response: %s", response.Status)

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling authentication service"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	// write the json response
	app.writeJSON(w, http.StatusAccepted, payload)
}
