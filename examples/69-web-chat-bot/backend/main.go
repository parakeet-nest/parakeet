package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/history"
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

	ollamaUrl := os.Getenv("OLLAMA_BASE_URL")
	if ollamaUrl == "" {
		ollamaUrl = "http://localhost:11434"
	}

	model := os.Getenv("LLM_CHAT")
	if model == "" {
		model = "deepseek-r1:1.5b"
	}

	var httpPort = os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "5050"
	}

	fmt.Println("🌍", ollamaUrl, "📕", model)

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.5,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 2.2,
		option.TopK:          10,
		option.TopP:          0.5,
	})

	systemInstructions := `You are a useful AI agent, your name is Bob`

	conversation := history.MemoryMessages{
		Messages: make(map[string]llm.MessageRecord),
	}

	mux := http.NewServeMux()
	shouldIStopTheCompletion := false

	mux.HandleFunc("POST /chat", func(response http.ResponseWriter, request *http.Request) {
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

		userMessage := data["message"]
		previousMessages, _ := conversation.GetAllMessages()

		// (Re)Create the conversation
		conversationMessages := []llm.Message{}
		// instruction
		conversationMessages = append(conversationMessages, llm.Message{Role: "system", Content: systemInstructions})
		// history
		conversationMessages = append(conversationMessages, previousMessages...)
		// last question
		conversationMessages = append(conversationMessages, llm.Message{Role: "user", Content: userMessage})

		fmt.Println("🅰:", conversationMessages)

		query := llm.Query{
			Model:    model,
			Messages: conversationMessages,
			Options:  options,
		}
		/*
			query := llm.Query{
				Model: model,
				Messages: []llm.Message{
					{Role: "system", Content: systemInstructions},
					{Role: "user", Content: userMessage},
				},
				Options: options,
			}
		*/

		answer, err := completion.ChatStream(ollamaUrl, query,
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

		conversation.SaveMessage(uuid.New().String(), llm.Message{
			Role:    "user",
			Content: userMessage,
		})
		conversation.SaveMessage(uuid.New().String(), llm.Message{
			Role:    "system",
			Content: answer.Message.Content,
		})

	})

	// Cancel/Stop the generation of the completion
	mux.HandleFunc("DELETE /cancel", func(response http.ResponseWriter, request *http.Request) {
		shouldIStopTheCompletion = true
		response.Write([]byte("🚫 Cancelling request..."))
	})

	var errListening error
	log.Println("🌍 http server is listening on: " + httpPort)
	errListening = http.ListenAndServe(":"+httpPort, mux)

	log.Fatal(errListening)

}
