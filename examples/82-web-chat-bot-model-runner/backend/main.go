package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/enums/provider"
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

	modelRunnerURL := os.Getenv("MODEL_RUNNER_BASE_URL") + "/engines/llama.cpp/v1"

	model := os.Getenv("LLM_CHAT")

	fmt.Println("modelRunnerURL:", modelRunnerURL)
	fmt.Println("model:", model)

	var httpPort = os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "5050"
	}

	fmt.Println("🌍", modelRunnerURL, "📕", model)

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.5,
		option.RepeatPenalty: 2.2,
	})

	systemInstructions := `You are a useful AI agent, your name is Bob`

	conversation := history.MemoryMessages{
		Messages: make(map[string]llm.MessageRecord),
	}

	mux := http.NewServeMux()
	shouldIStopTheCompletion := false

	conversationLength := 6

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

		/*
		// (Re)Create the conversation
		conversationMessages := []llm.Message{}
		// instruction
		conversationMessages = append(conversationMessages, llm.Message{Role: "system", Content: systemInstructions})
		// history
		conversationMessages = append(conversationMessages, previousMessages...)
		// last question
		conversationMessages = append(conversationMessages, llm.Message{Role: "user", Content: userMessage})
		*/

		// (Re)Create the conversation
		conversationMessages := llm.SetOfMessages(
			llm.Message{Role: "system", Content: systemInstructions},
			previousMessages,
			llm.Message{Role: "user", Content: userMessage},
		)

		//? 📝 Print the previous messages
		fmt.Println("👋 previousMessages:")
		for _, message := range previousMessages {
			fmt.Println(" - message:", message)
		}

		query := llm.Query{
			Model:    model,
			Messages: conversationMessages,
			Options:  options,
		}

		answer, err := completion.ChatStream(modelRunnerURL, query,
			func(answer llm.Answer) error {
				response.Write([]byte(answer.Message.Content))

				flusher.Flush()
				if !shouldIStopTheCompletion {
					return nil
				} else {
					return errors.New("🚫 Cancelling request")
				}
			}, provider.DockerModelRunner)

		if err != nil {
			shouldIStopTheCompletion = false
			response.Write([]byte("bye: " + err.Error()))
		}

		conversation.SaveMessage("", llm.Message{
			Role:    "user",
			Content: userMessage,
		})

		conversation.SaveMessage("", llm.Message{
			Role:    "assistant",
			Content: answer.Message.Content,
		})

		conversation.KeepLastN(conversationLength)

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
