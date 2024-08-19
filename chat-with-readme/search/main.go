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
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	//ollamaUrl := "http://bob.local:11434"
	
	embeddingsModel := "all-minilm:33m"
	//embeddingsModel := "codestral"
	//embeddingsModel := "magicoder:latest"
	//embeddingsModel := "phi3:mini"

	//chatModel := "magicoder:latest"
	//chatModel := "codestral"
	//chatModel := "deepseek-coder:6.7b"
	chatModel := "phi3:mini" 
	//chatModel := "llama3" 
	//chatModel := "granite-code:3b"
	//chatModel := "gemma2:2b"


	store := embeddings.BboltVectorStore{}
	store.Initialize("../embeddings.db")

	systemContent := `You are an expert with the Parakeet library.
	Please make friendly answer for the noobs. Use only the provided document to answer.
	Add source code examples if you can, only in Golang.`

	// ✋ it's important to explain the LLM that it must use only the context and doc
	// otherwise it will use first its knowledge and then the answer
	// will be less accurate

	// Question for the Chat system
	//userContent := `[Brief] How to create a stream completion with Parakeet?`
	//userContent := `How to create a simple completion with Parakeet?`
	
	//userContent := `Explain how to create a stream chat completion with Parakeet?`

	userContent := `Verbose mode`

	//userContent := `Explain how can I generate JSON output`

	// Create an embedding from the user question
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
	//fmt.Println(embeddingFromQuestion)

	fmt.Println("🔎 searching for similarity...")

	//fmt.Println(store.GetAll())

	similarities, err := store.SearchSimilarities(embeddingFromQuestion, 0.2)

	//similarities, err := store.SearchTopNSimilarities(embeddingFromQuestion, 0.3, 3)

	if err != nil {
		fmt.Println("😡", err)
	}

	documentsContent := embeddings.GenerateContextFromSimilarities(similarities)

	fmt.Println("🎉 similarities", len(similarities))

	query := llm.Query{
		Model: chatModel,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "system", Content: documentsContent},
			{Role: "user", Content: userContent},
		},
		Options: llm.Options{
			Temperature: 0.8,
			RepeatLastN: 2,
		},
		//Stream: false,
	}

	fmt.Println("")
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

}
