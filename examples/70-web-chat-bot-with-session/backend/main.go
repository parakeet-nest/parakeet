package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

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

var m sync.Mutex
var messagesCounters = make(map[string]int)

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

	messagesCounters := map[string]int{}

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

		fmt.Println("📝 posted data:", data)

		userMessage := data["message"]
		sessionId := data["sessionId"]

		//? 💌 Get all messages from the conversation filtered by the session id
		//previousMessages, _ := conversation.GetAllMessages()
		previousMessages, _ := conversation.GetAllMessagesOfSession(sessionId)

		//? 📝 Print the previous messages
		fmt.Println("👋 sessionId", sessionId, "previousMessages:")
		for _, message := range previousMessages {
			fmt.Println(" - message:", message)
		}

		// (Re)Create the conversation
		conversationMessages := []llm.Message{}
		// instruction
		conversationMessages = append(conversationMessages, llm.Message{Role: "system", Content: systemInstructions})
		// history
		conversationMessages = append(conversationMessages, previousMessages...)
		// last question
		conversationMessages = append(conversationMessages, llm.Message{Role: "user", Content: userMessage})

		query := llm.Query{
			Model:    model,
			Messages: conversationMessages,
			Options:  options,
		}

		answer, err := completion.ChatStream(ollamaUrl, query,
			func(answer llm.Answer) error {
				//log.Println("📝:", answer.Message.Content)
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

		// Is it useful or not?
		m.Lock()
		defer m.Unlock()
		//! I use a counter for the id of the message, then I can create an ordered list of messages

		conversation.SaveMessageWithSession(sessionId, &messagesCounters, llm.Message{
			Role:    "user",
			Content: userMessage,
		})
		//* Remove the top(first) message of conversation of maxMessages(6) messages
		conversation.RemoveTopMessageOfSession(sessionId, &messagesCounters, 6)

		conversation.SaveMessageWithSession(sessionId, &messagesCounters, llm.Message{
			Role:    "assistant",
			Content: answer.Message.Content,
		})
		conversation.RemoveTopMessageOfSession(sessionId, &messagesCounters, 6)

	})

	mux.HandleFunc("POST /clear-history", func(response http.ResponseWriter, request *http.Request) {
		// TODO: Clear all messages from the conversation
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
