package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	embeddingsModel := "all-minilm:33m" // This model is for the embeddings of the documents
	//smallChatModel := "qwen2:0.5b"      // This model is for the chat completion
	//smallChatModel := "qwen:0.5b"      // This model is for the chat completion
	smallChatModel := "tinydolphin"      // This model is for the chat completion

	/*
	run a container from the hello-world docker image
	get the list of the running containers
	get the list of the images
	how to build a docker image

	*/


	store := embeddings.BboltVectorStore{}
	store.Initialize("../embeddings.db")

	for {
		question := input()
		if question == "quit" {
			break
		}
		getCompletion(question, ollamaUrl, embeddingsModel, smallChatModel, store, 0.6)

	}

}

func input() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("🐋🤖 ask me something> ")
	question, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln("😡:", err)
		return ""
	}

	question = strings.TrimSpace(question)
	return question
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

	systemContent := `translate this sentence in docker command, using only the provided context`

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
