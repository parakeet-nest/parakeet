package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/sea-monkeys/daphnia"
)

func TestQuestionContextChunks(t *testing.T) {

	generateContextFromSimilarities := func (similarities []daphnia.VectorRecord) string {
		documentsContent := "<context>\n"
		for _, similarity := range similarities {
			documentsContent += fmt.Sprintf("<doc>%s</doc>\n", similarity.Prompt)
		}
		documentsContent += "</context>"
		return documentsContent
	}


	ollamaUrl := os.Getenv("OLLAMA_URL")
	if ollamaUrl == "" {
		ollamaUrl = "http://localhost:11434"
	}
	embeddingsModel := "mxbai-embed-large:latest"
	chatModel := "qwen2.5:1.5b"

	/*
		options := llm.SetOptions(map[string]interface{}{
			option.Temperature: 0.0,
		})
	*/
	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.0,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 3.0,
		option.TopK:          10,
		option.TopP:          0.5,
	})

	// Initialize the vector store
	vectorStore := daphnia.VectorStore{}
	vectorStore.Initialize("with-context.gob")

	question := "Explain the biological compatibility of the Human species?"

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

	embeddingForDaphnia := daphnia.VectorRecord{
		Prompt:    embeddingFromQuestion.Prompt,
		Embedding: embeddingFromQuestion.Embedding,
	}

	similarities, _ := vectorStore.SearchTopNSimilarities(embeddingForDaphnia, 0.65, 10)

	for _, similarity := range similarities {
		fmt.Println()
		fmt.Println("Cosine distance:", similarity.CosineDistance)

		fmt.Println(similarity.Prompt)
	}

	if len(similarities) == 0 {
		t.Errorf("No similarities found")
	} else {
		fmt.Println("🎉 number of similarities:", len(similarities))
	}

	documentsContent := generateContextFromSimilarities(similarities)

	messages := []llm.Message{
		{Role: "system", Content: "You are a usefull AI agent, expert with Heroic Fantasy and Science Fiction. Use only the provides content to answer."},
		{Role: "system", Content: documentsContent},
		{Role: "user", Content: question},
	}

	queryChat := llm.Query{
		Model:    chatModel,
		Messages: messages,
		Options:  options,
	}

	fmt.Println()
	fmt.Println("🤖 answer:")

	// Answer the question
	_, errCompletion := completion.ChatStream(ollamaUrl, queryChat,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if errCompletion != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println()

}


