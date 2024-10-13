package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/ui"
	"github.com/parakeet-nest/parakeet/ui/colors"

	"fmt"
	"log"
)

func ChatWithCharacter(instructions, description, ollamaUrl, model string) {

	systemContent := instructions

	contextContext := description

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.5,
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
				{Role: "system", Content: systemContent},
				{Role: "system", Content: contextContext},
				{Role: "user", Content: question},
			},
			Options:          options,
		}

		fmt.Println()
		ui.Println(colors.Magenta, "🤖 answer:")

		// Answer the question
		_, err := completion.ChatStream(ollamaUrl, queryChat,
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
		model = "nemotron-mini"
	}
	// nemotron-mini 🤩
	// qwen2:1.5b 🙂
	// gemma2:2b 🙂
	// dolphin-phi:2.7b 🙂

	// some questions to ask:
	// what is your name?
	// give me the list without detail of your qualities
	// where are you from?
	// where are you located?
	// where are you living?
	// where did you grow up?
	// who is your best friend?
	// who is your worst enemy?
	// give me the list without detail of all your friends
	// give me the list without detail of all your enemies

	instructions, err := os.ReadFile("instructions.md")
	if err != nil {
		log.Fatal("😡:", err)
	}
	description, err := os.ReadFile("description.md")
	if err != nil {
		log.Fatal("😡:", err)
	}

	ChatWithCharacter(
		string(instructions),
		string(description),
		ollamaUrl,
		model,
	)
}
