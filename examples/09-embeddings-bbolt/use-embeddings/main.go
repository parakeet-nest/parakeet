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
	var embeddingsModel = "all-minilm" // This model is for the embeddings of the documents
	var smallChatModel = "qwen:0.5b"   // This model is for the chat completion

	store := embeddings.BboltVectorStore{}
	store.Initialize("../embeddings.db")

	// Question for the Chat system
	userContent := `Who is Philippe Charrière and what spaceship does he work on?`
	//userContent := `What is the nickname of Philippe Charrière?`

	systemContent := `You are an AI assistant. Your name is Seven. 
		Some people are calling you Seven of Nine.
		You are an expert in Star Trek.
		All questions are about Star Trek.
		Using the provided context, answer the user's question
		to the best of your ability using only the resources provided.`

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

	similarity, _ := store.SearchMaxSimilarity(embeddingFromQuestion)

	//similarities, _ := store.SearchTopNSimilarities(embeddingFromQuestion, 0.3, 1)
	//similarities, _ := store.SearchSimilarities(embeddingFromQuestion, 0.3)
	//similarity := similarities[0]

	fmt.Println("🎉 similarity", similarity)

	documentsContent := `<context><doc>` + similarity.Prompt + `</doc></context>`

	query := llm.Query{
		Model: smallChatModel,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "system", Content: documentsContent},
			{Role: "user", Content: userContent},
		},
		Options: llm.Options{
			Temperature: 0.4,
			RepeatLastN: 2,
		},
		Stream: false,
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
