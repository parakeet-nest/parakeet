package main

import (
	"fmt"
	"log"


	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	var embeddingsModel = "all-minilm:33m" // This model is for the embeddings of the documents
	var smallChatModel = "qwen2:0.5b"      // This model is for the chat completion

	store := embeddings.BboltVectorStore{}
	store.Initialize("../embeddings.db")

	getCompletion("List all containers running prior to 3e3o3ad9a0b2e.", ollamaUrl, embeddingsModel, smallChatModel, store, 0.3)

	getCompletion("How to get the list of the docker images on my computer?", ollamaUrl, embeddingsModel, smallChatModel, store, 0.3)

}

func getContentFromSimilarities(userContent, ollamaUrl, embeddingsModel string, store embeddings.BboltVectorStore, limit float64) string {

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

	similarities, _ := store.SearchSimilarities(embeddingFromQuestion, limit)

	fmt.Println("🎉 similarities:", len(similarities))

	documentsContent := embeddings.GenerateContentFromSimilarities(similarities)

	return documentsContent

}

func getCompletion(userContent, ollamaUrl, embeddingsModel, smallChatModel string, store embeddings.BboltVectorStore, limit float64) {

	systemContent := `translate this sentence in docker command`

	documentsContent := getContentFromSimilarities(userContent, ollamaUrl, embeddingsModel, store, 0.3)

	query := llm.Query{
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
	_, err := completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}

	fmt.Println()

}
