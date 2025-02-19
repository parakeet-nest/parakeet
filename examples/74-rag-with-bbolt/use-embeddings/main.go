package main

import (
	"fmt"
	"log"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/enums/option"

)

func main() {
	ollamaUrl := "http://localhost:11434"
	var embeddingsModel = "all-minilm:33m" // This model is for the embeddings of the documents
	var smallChatModel = "qwen2.5:7b"      // This model is for the chat completion

	store := embeddings.BboltVectorStore{}
	//err := store.Initialize("../embeddings.readme.db")
	//err := store.Initialize("../embeddings.db")
	err := store.Initialize("../embeddings.updated.db")


	//vectors, _ := store.GetAll()
	//fmt.Println(vectors)

	if err != nil {
		log.Fatalln("😡:", err)
	}

	systemContent := `You are a Golang expert and know very well the extism go SDK`

	//userContent := `What is Parakeet`
	userContent := `How to call a function of a wasm module in golang with the extism-go sdk?`

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
	similarities, _ := store.SearchTopNSimilarities(embeddingFromQuestion, 0.7, 2)
	//similarities, _ := store.SearchTopNSimilarities(embeddingFromQuestion, 0.3, 3)


	for _, similarity := range similarities {
		// Do something with the similarity
		fmt.Println("Similarity:", similarity.SimpleMetaData)
	}

	fmt.Println("🎉 number of similarities:", len(similarities))

	documentsContent := embeddings.GenerateContextFromSimilarities(similarities)

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.0,
		option.RepeatLastN: 2,
		option.RepeatPenalty: 2.2,
		option.TopK: 10,
		option.TopP: 0.5,
	})

	query := llm.Query{
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
