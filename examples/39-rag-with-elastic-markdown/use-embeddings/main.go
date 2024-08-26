package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡:", err)
	}

	ollamaUrl := "http://localhost:11434"
	embeddingsModel := "all-minilm:33m" // This model is for the embeddings of the documents
	smallChatModel := "qwen2:0.5b"      // This model is for the chat completion

	cert, _ := os.ReadFile(os.Getenv("ELASTIC_CERT_PATH"))

	elasticStore := embeddings.ElasticSearchStore{}
	err = elasticStore.Initialize(
		[]string{
			os.Getenv("ELASTIC_ADDRESS"),
		},
		os.Getenv("ELASTIC_USERNAME"),
		os.Getenv("ELASTIC_PASSWORD"),
		cert,
		"golang-index",
	)
	if err != nil {
		log.Fatalln("😡:", err)
	}

	//userContent := `Who are the monsters of Chronicles of Aethelgard?`
	//userContent := `Tell me more about Keegorg`
	userContent := `What changes to the archive/tar library happened in Go 1.23`

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

	queryChat := llm.Query{
		Model: smallChatModel,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "system", Content: documentsContent},
			{Role: "user", Content: userContent},
		},
		Options: llm.Options{
			Temperature:   0.0,
			RepeatLastN:   2,
			RepeatPenalty: 3.0,
			TopK:          10,
			TopP:          0.5,
		},
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
