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
	err := store.Initialize("../embeddings.db")

	//vectors, _ := store.GetAll()
	//fmt.Println(vectors)

	if err != nil {
		log.Fatalln("😡:", err)
	}

	systemContent := `You are a Golang expert and know very well the GenAI library Parakeet`

	userContent := `What is Parakeet`

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

	//similarities, _ := store.SearchSimilarities(embeddingFromQuestion, 0.3)
	similarities, _ := store.SearchTopNSimilarities(embeddingFromQuestion, 0.3, 2)

	for _, similarity := range similarities {
		// Do something with the similarity
		fmt.Println("Similarity:", similarity.MetaData)
	}

	fmt.Println("🎉 number of similarities:", len(similarities))

	documentsContent := embeddings.GenerateContextFromSimilarities(similarities)

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
	_, err = completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}

	fmt.Println()
}
