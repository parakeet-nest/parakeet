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
	"github.com/parakeet-nest/parakeet/gear"
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

// USed for the history of messages
var m sync.Mutex
var messagesCounters = make(map[string]int)

func main() {
	httpPort := gear.GetEnvString("HTTP_PORT", "5050")

	conversation := history.MemoryMessages{
		Messages: make(map[string]llm.MessageRecord),
	}
	//messagesCounters := map[string]int{}

	var directoryTreePath = gear.GetEnvString("DIRECTORY_TREE_PATH", "../data/tree.txt")

	directoryTree, err := os.ReadFile(directoryTreePath)
	if err != nil {
		log.Fatalln("😡:", err)
	}

	var contentPath = gear.GetEnvString("CONTENT_PATH", "../data/content.txt")

	contentFiles, err := os.ReadFile(contentPath)
	if err != nil {
		log.Fatalln("😡:", err)
	}

	systemInstructionsPath := gear.GetEnvString("SYSTEM_INSTRUCTIONS_PATH", "../instructions/parakeet.instructions.md")

	systemInstructions, err := os.ReadFile(systemInstructionsPath)
	if err != nil {
		log.Fatalln("😡:", err)
	}

	fmt.Println("🤖📝 system instructions:", string(systemInstructions))

	ollamaUrl := gear.GetEnvString("OLLAMA_BASE_URL", "http://localhost:11434")

	model := gear.GetEnvString("LLM_CHAT", "deepseek-r1:1.5b")

	// Options
	temperature := gear.GetEnvFloat("OPTION_TEMPERATURE", 0.5)
	repeatLastN := gear.GetEnvInt("OPTION_REPEAT_LAST_N", 2)
	repeatPenalty := gear.GetEnvFloat("OPTION_REPEAT_PENALTY", 2.2)
	topK := gear.GetEnvInt("OPTION_TOP_K", 10)
	topP := gear.GetEnvFloat("OPTION_TOP_P", 0.5)
	minP := gear.GetEnvFloat("OPTION_MIN_P", 0.1)
	microstat := gear.GetEnvInt("OPTION_MIROSTAT", 1)
	microstatTau := gear.GetEnvFloat("OPTION_MIROSTAT_TAU", 3.0)
	microstatEta := gear.GetEnvFloat("OPTION_MIROSTAT_ETA", 0.1)



	fmt.Println("🌍", ollamaUrl, "📕", model)

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   temperature,
		option.RepeatLastN:   repeatLastN,
		option.RepeatPenalty: repeatPenalty,
		option.TopK:          topK,
		option.TopP:          topP,
		option.MinP:          minP,
		option.Mirostat:      microstat,
		option.MirostatTau:   microstatTau,
		option.MirostatEta:   microstatEta,
	})

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

		fmt.Println("📝 posted data:", data)

		userMessage := data["message"]
		sessionId := data["sessionId"]
		fmt.Println("📝 sessionId:", sessionId)

		//? History of messages
		previousMessages, _ := conversation.GetAllMessagesOfSession(sessionId)
		//? End of history of messages

		// (Re)Create the conversation
		conversationMessages := []llm.Message{}

		// history
		conversationMessages = append(conversationMessages, previousMessages...)

		// instruction
		conversationMessages = append(conversationMessages, llm.Message{Role: "system", Content: string(systemInstructions) + "\n"})
		// history
		//conversationMessages = append(conversationMessages, previousMessages...)
		// Repository tree
		conversationMessages = append(conversationMessages, llm.Message{Role: "system", Content: "REPOSITORY:\n" + string(directoryTree)})
		// Source code
		conversationMessages = append(conversationMessages, llm.Message{Role: "system", Content: "CONTENT:\n" + string(contentFiles)})
		// last question
		conversationMessages = append(conversationMessages, llm.Message{Role: "user", Content: userMessage})

		// Prompt construction
		/*
			messages := []llm.Message{
				{Role: "system", Content: string(systemInstructions) + "\n"},
				{Role: "system", Content: "REPOSITORY:\n" + string(directoryTree)},
				{Role: "system", Content: "SOURCE CODE:\n" + repositoryContent},
				{Role: "user", Content: userMessage},
			}
		*/

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
		//* Remove the top(first) message of conversation of maxMessages(3) messages
		conversation.RemoveTopMessageOfSession(sessionId, &messagesCounters, 3)

		conversation.SaveMessageWithSession(sessionId, &messagesCounters, llm.Message{
			Role:    "assistant",
			Content: answer.Message.Content,
		})
		conversation.RemoveTopMessageOfSession(sessionId, &messagesCounters, 3)

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
