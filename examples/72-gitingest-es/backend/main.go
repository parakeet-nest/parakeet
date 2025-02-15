package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/embeddings"
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

	var contentPath = os.Getenv("CONTENT_PATH")
	if contentPath == "" {
		contentPath = "../data/tree.txt"
	}

	// open ../data/tree.txt
	// read the content of tree.txt
	directoryTree, err := os.ReadFile(contentPath)
	if err != nil {
		log.Fatal(err)
	}

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

	embeddingsModel := os.Getenv("LLM_EMBEDDINGS")
	if embeddingsModel == "" {
		embeddingsModel = "mxbai-embed-large"
	}

	maxSimilaritiesEnv := os.Getenv("MAX_SIMILARITIES")
	if maxSimilaritiesEnv == "" {
		maxSimilaritiesEnv = "5"
	}
	maxSimilarities , err := strconv.Atoi(maxSimilaritiesEnv)
	if err != nil {
		maxSimilarities = 5
	}

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
		option.Temperature:   0.5,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 2.2,
		option.TopK:          10,
		option.TopP:          0.5,
	})

	systemInstructions := os.Getenv("SYSTEM_INSTRUCTIONS")
	if systemInstructions == "" {
		systemInstructions = `You are a useful AI agent, your name is Bob`
	}

	fmt.Println("🤖📝 system instructions:", systemInstructions)


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

		// Prompt construction
		messages := []llm.Message{
			{Role: "system", Content: string(systemInstructions)+"\n"},
			{Role: "system", Content: "REPOSITORY:\n" + string(directoryTree)},
			{Role: "system", Content: "SOURCE CODE:\n" + repositoryContent},
			{Role: "user", Content: userMessage},
		}

		query := llm.Query{
			Model:    model,
			Messages: messages,
			Options:  options,
		}

		_, err = completion.ChatStream(ollamaUrl, query,
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
