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

	//smallChatModel := "qwen2:0.5b"
	smallChatModel := "gemma2:2b"
	embeddingsModel := "all-minilm:33m"


	systemContent := `instruction: 
	translate the user question in docker command using the given context.
	Stay brief.`

	store := embeddings.BboltVectorStore{}
	store.Initialize("../embeddings.db")

	for {
		question := input(smallChatModel)
		if question == "bye" {
			break
		}

		// Create an embedding from the question
		embeddingFromQuestion, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: question,
			},
			"question",
		)
		if err != nil {
			log.Fatalln("😡:", err)
		}
		fmt.Println("🔎 searching for similarity...")
		similarities, _ := store.SearchTopNSimilarities(embeddingFromQuestion, 0.4, 3)

		contextContent := embeddings.GenerateContextFromSimilarities(similarities)
		//fmt.Println(documentsContent)
		fmt.Println("🎉 similarities:", len(similarities))

		// Prepare the query
		query := llm.Query{
			Model: smallChatModel,
			Messages: []llm.Message{
				{Role: "system", Content: systemContent},
				{Role: "system", Content: contextContent},
				{Role: "user", Content: question},
			},
			Options: llm.Options{
				Temperature:   0.0,
				RepeatLastN:   2,
				RepeatPenalty: 3.0,
				TopK:          10,
				TopP:          0.5,
			},
		}

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
}

func input(smallChatModel string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("🐳 [%s] ask me something> ", smallChatModel)
	question, _ := reader.ReadString('\n')
	return strings.TrimSpace(question)
}
