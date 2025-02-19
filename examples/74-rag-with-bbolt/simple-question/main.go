package main

import (
	"fmt"
	"log"
	"os"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	var smallChatModel = "qwen2.5:7b" // This model is for the chat completion

	systemContent := `You are a Golang expert and know very well the extism go SDK. Use only the provided content to answer the question.`

	contentFiles, errFile := os.ReadFile("./content.md")
	if errFile != nil {
		log.Fatalln("😡:", errFile)
	}

	documentsContent := string(contentFiles)

	//userContent := `What is Parakeet`
	userContent := `How to call a function of a wasm module in golang with the extism-go sdk?`

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.0,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 2.2,
		option.TopK:          10,
		option.TopP:          0.5,
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
