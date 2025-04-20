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
	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/embeddings"
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

	systemInstructionsPath := gear.GetEnvString("SYSTEM_INSTRUCTIONS_PATH", "../instructions/parakeet.instructions.md")

	systemInstructions, err := os.ReadFile(systemInstructionsPath)
	if err != nil {
		log.Fatalln("😡:", err)
	}

	fmt.Println("🤖📝 system instructions:", string(systemInstructions))

	ollamaUrl := gear.GetEnvString("OLLAMA_BASE_URL", "http://localhost:11434")

	model := gear.GetEnvString("LLM_CHAT", "qwen2.5:3b")
	embeddingsModel := gear.GetEnvString("LLM_EMBEDDINGS", "mxbai-embed-large")

	maxSimilarities := gear.GetEnvInt("MAX_SIMILARITIES", 5)

	historyMessages := gear.GetEnvInt("HISTORY_MESSAGES", 3)

	// Options
	temperature := gear.GetEnvFloat("OPTION_TEMPERATURE", 0.5)
	repeatLastN := gear.GetEnvInt("OPTION_REPEAT_LAST_N", 2)
	repeatPenalty := gear.GetEnvFloat("OPTION_REPEAT_PENALTY", 2.2)
	topK := gear.GetEnvInt("OPTION_TOP_K", 10)
	topP := gear.GetEnvFloat("OPTION_TOP_P", 0.5)
	numCtx := gear.GetEnvInt("OPTION_NUM_CTX", 4096)

	// Initialize the Elasticsearch store
	elasticStore := embeddings.ElasticsearchStore{}
	err = elasticStore.Initialize(
		[]string{
			os.Getenv("ELASTICSEARCH_HOSTS"),
		},
		os.Getenv("ELASTICSEARCH_USERNAME"),
		os.Getenv("ELASTICSEARCH_PASSWORD"),
		nil,
		os.Getenv("ELASTICSEARCH_INDEX"),
	)
	if err != nil {
		log.Fatalln("😡:", err)
	}

	fmt.Println("🌍", ollamaUrl, "📕", model, "🌐", embeddingsModel)

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   temperature,
		option.RepeatLastN:   repeatLastN,
		option.RepeatPenalty: repeatPenalty,
		option.TopK:          topK,
		option.TopP:          topP,
		option.NumCtx:        numCtx,
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

		//! Similarity search
		// Create an embedding from the question
		embeddingFromQuestion, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: userMessage,
			},
			"question",
		)
		if err != nil {
			log.Fatalln("😡:", err)
		}
		fmt.Println("🔎 searching for similarity...")

		similarities, err := elasticStore.SearchTopNSimilarities(embeddingFromQuestion, maxSimilarities)

		for _, similarity := range similarities {
			fmt.Println("📝 doc:", similarity.Id, "score:", similarity.Score)
			fmt.Println("--------------------------------------------------")
			fmt.Println("📝 metadata:", similarity.Prompt)
			fmt.Println("--------------------------------------------------")
		}

		if err != nil {
			log.Fatalln("😡:", err)
		}

		repositoryContent := embeddings.GenerateContentFromSimilarities(similarities)

		//! End of similarity search

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
		//conversationMessages = append(conversationMessages, llm.Message{Role: "system", Content: "SOURCE CODE:\n" + repositoryContent})
		conversationMessages = append(conversationMessages, llm.Message{Role: "system", Content:  repositoryContent})

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

		// Estimate the number of tokens
		concatenatedMessages:= ""
		for _, message := range conversationMessages {
			//fmt.Println("📝", message.Content)
			concatenatedMessages += message.Content + "\n"
		}
		fmt.Println("================================================")
		estimatedTokens := content.EstimateGPTTokens(concatenatedMessages)
		fmt.Println("🧩 estimated tokens:", estimatedTokens)
		fmt.Println("================================================")

		if numCtx < estimatedTokens {
			fmt.Println("🔥 numCtx is less than estimated tokens")
			options.NumCtx = estimatedTokens + 100
		} else {
			options.NumCtx = numCtx
		}

		
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

		conversation.SaveMessageWithSession(sessionId, "", llm.Message{
			Role:    "user",
			Content: userMessage,
		})

		conversation.SaveMessageWithSession(sessionId, "", llm.Message{
			Role:    "assistant",
			Content: answer.Message.Content,
		})
		conversation.KeepLastNOfSession(sessionId, historyMessages)

		//* Create an embedding from the user message if it starts with "LEARN:"
		/*
		if userMessage[:6] == "LEARN:" {
			embedding, err := embeddings.CreateEmbedding(
				ollamaUrl,
				llm.Query4Embedding{
					Model:  embeddingsModel,
					Prompt: userMessage[6:],
				},
				"learn",
			)
			if err != nil {
				log.Fatalln("😡:", err)
			}

			if _, err = elasticStore.Save(embedding); err != nil {
				log.Fatalln("😡:", err)
			}

			fmt.Println("🎉 Document", embedding.Id, "indexed successfully")
		}
		*/

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
