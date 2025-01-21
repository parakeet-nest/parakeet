package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡:", err)
	}

	ollamaUrl := "http://localhost:11434"
	//embeddingsModel := "all-minilm:33m" // This model is for the embeddings of the documents
	//embeddingsModel := "nomic-embed-text"
	embeddingsModel := "mxbai-embed-large"
	//smallChatModel := "gemma2:2b"      // This model is for the chat completion
	//smallChatModel := "phi3:mini"      // This model is for the chat completion
	smallChatModel := "llama3.1:8b"

	elasticStore := embeddings.ElasticsearchStore{}
	err = elasticStore.Initialize(
		[]string{
			os.Getenv("ELASTIC_ADDRESS"),
		},
		os.Getenv("ELASTIC_USERNAME"),
		os.Getenv("ELASTIC_PASSWORD"),
		nil,
		"hierarchy-mxbai-golang-index",
	)
	if err != nil {
		log.Fatalln("😡:", err)
	}

	//userContent := `[Brief] What's new with TLS client?`
	//userContent := `Tell me more about the new structs package`
	userContent := `Summarize what's new with benchmarks in 3 bullet points. Be succinct`

	// Create an embedding from the question
	embeddingFromQuestion, err := embeddings.CreateEmbedding(
		ollamaUrl,
		llm.Query4Embedding{
			Model:  embeddingsModel,
			Prompt: userContent,
		},
		"question",
	)
	if err != nil {
		log.Fatalln("😡:", err)
	}
	fmt.Println("🔎 searching for similarity...")

	similarities, err := elasticStore.SearchTopNSimilarities(embeddingFromQuestion, 3)

	for _, similarity := range similarities {
		fmt.Println("📝 doc:", similarity.Id, "score:", similarity.Score)
	}

	if err != nil {
		log.Fatalln("😡:", err)
	}

	documentsContent := embeddings.GenerateContentFromSimilarities(similarities)

	systemContent := `You are a Golang expert.
	Using only the below provided context, answer the user's question
	to the best of your ability using only the resources provided.
	`

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.0,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 3.0,
		option.TopK:          10,
		option.TopP:          0.5,
	})

	queryChat := llm.Query{
		Model: smallChatModel,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "system", Content: documentsContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
	}

	fmt.Println()
	fmt.Println("🤖 answer:")

	// Answer the question
	_, err = completion.ChatStream(ollamaUrl, queryChat,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})
	if err != nil {
		log.Fatal("😡:", err)
	}

	fmt.Println()
}
