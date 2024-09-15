package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/llm"
)

/*
GetBytesBody returns the body of an HTTP request as a []byte.
  - It takes a pointer to an http.Request as a parameter.
  - It returns a []byte.
*/
func GetBytesBody(request *http.Request) []byte {
	body := make([]byte, request.ContentLength)
	request.Body.Read(body)
	return body
}

func main() {

	var ollamaUrl = os.Getenv("OLLAMA_BASE_URL")
	if ollamaUrl == "" {
		ollamaUrl = "http://localhost:11434"
	}

	var httpPort = os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	mux := http.NewServeMux()

	fileServerHtml := http.FileServer(http.Dir("public"))
	mux.Handle("/", fileServerHtml)

	shouldIStopTheCompletion := false

	mux.HandleFunc("GET /api/models", func(response http.ResponseWriter, request *http.Request) {
		modelsList, statusCode, err := llm.GetModelsList(ollamaUrl)
		if err != nil {
			response.Write([]byte("😡 Error: " + err.Error()))
		}

		jsonList, err := json.Marshal(&modelsList.Models)
		if err != nil {
			response.Write([]byte("😡 Error: " + err.Error()))
		}
		if statusCode != 200 {
			response.Write([]byte("😡 Error: " + http.StatusText(statusCode)))
		}

		response.Header().Add("Content-Type", "application/json; charset=utf-8")
		response.Write(jsonList)
	})

	mux.HandleFunc("POST /api/settings/chat", func(response http.ResponseWriter, request *http.Request) {
		// add a flusher
		flusher, ok := response.(http.Flusher)
		if !ok {
			response.Write([]byte("😡 Error: expected http.ResponseWriter to be an http.Flusher"))
		}
		body := GetBytesBody(request)
		// unmarshal the json data
		var data map[string]string

		err := json.Unmarshal(body, &data)
		if err != nil {
			response.Write([]byte("😡 Error: " + err.Error()))
		}

		model := data["model"]
		systemContent := data["system"]
		userContent := data["user"]

		temperatureStr := data["temperature"]
		repeatLastNStr := data["repeatLastN"]
		repeatPenaltyStr := data["repeatPenalty"]

		repeatPenalty, err := strconv.ParseFloat(repeatPenaltyStr, 64)
		if err != nil {
			response.Write([]byte("😡 Error: invalid repeatPenalty value"))
			return
		}

		temperature, err := strconv.ParseFloat(temperatureStr, 64)
		if err != nil {
			response.Write([]byte("😡 Error: invalid temperature value"))
			return
		}
		if temperature == 0 {
			temperature = 0.0
		}

		repeatLastN, err := strconv.Atoi(repeatLastNStr)
		if err != nil {
			response.Write([]byte("😡 Error: invalid repeatLastN value"))
			return
		}

		options := llm.SetOptions(map[string]interface{}{
			option.Temperature:   temperature,
			option.RepeatLastN:   repeatLastN,
			option.RepeatPenalty: repeatPenalty,
			option.Verbose:       true,
		})

		query := llm.Query{
			Model: model,
			Messages: []llm.Message{
				{Role: "system", Content: systemContent},
				{Role: "user", Content: userContent},
			},
			Options: options,
		}

		_, err = completion.ChatStream(ollamaUrl, query,
			func(answer llm.Answer) error {
				log.Println("📝:", answer.Message.Content)
				response.Write([]byte(answer.Message.Content))

				flusher.Flush()
				if !shouldIStopTheCompletion {
					return nil
				} else {
					return errors.New("🚫 Cancelling request")
				}
			})

		if err != nil {
			shouldIStopTheCompletion = false
			response.Write([]byte("bye: " + err.Error()))
		}

	})

	// Cancel/Stop the generation of the completion
	mux.HandleFunc("DELETE /api/completion/cancel", func(response http.ResponseWriter, request *http.Request) {
		shouldIStopTheCompletion = true
		response.Write([]byte("🚫 Cancelling request..."))
	})

	var errListening error
	log.Println("🌍 http server is listening on: " + httpPort)
	errListening = http.ListenAndServe(":"+httpPort, mux)

	log.Fatal(errListening)
}
