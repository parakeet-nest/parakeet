package main

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/ui"
	"github.com/parakeet-nest/parakeet/ui/colors"

	"fmt"
	"log"
)

func ChatWithCharacter(ollamaUrl, model string) {

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.0,
		option.RepeatLastN:   3,
		option.RepeatPenalty: 2.0,
		option.TopK:          10,
		option.TopP:          0.5,
	})

	for {

		question, _ := ui.Input(colors.Cyan, fmt.Sprintf("🤖 [%s] ask me something> ", model))

		if question == "bye" {
			break
		}

		queryChat := llm.Query{
			Model: model,
			Messages: []llm.Message{
				{Role: "user", Content: question},
			},
			Options: options,
			Stream:  false,
			//Format:  "json",
		}

		fmt.Println()
		ui.Println(colors.Magenta, "🤖 answer:")

		// Answer the question
		answer, err := completion.Chat(ollamaUrl, queryChat)
		if err != nil {
			log.Fatal("😡:", err)
		}

		if strings.HasPrefix(answer.Message.Content, "unsafe") {
			ui.Println(colors.Red, "😡", answer.Message.Content)
		} else {
			ui.Println(colors.Green, "🙂", answer.Message.Content)
		}



		fmt.Println()
	}

}

func main() {
	// create a `.env` file with the following content:
	// OLLAMA_URL=your_ollama_url

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡:", err)
	}

	ollamaUrl := os.Getenv("OLLAMA_URL")
	if ollamaUrl == "" {
		ollamaUrl = "http://localhost:11434"
	}

	model := os.Getenv("MODEL")
	if model == "" {
		model = "llama-guard3:1b"
	}

	ChatWithCharacter(
		ollamaUrl,
		model,
	)
}
