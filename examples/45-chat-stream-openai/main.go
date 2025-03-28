package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/provider"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
)

func main() {
	// create a `.env` file with the following content:
	// OPENAI_API_KEY=your_openai_api_key
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡:", err)
	}

	openAIURL := "https://api.openai.com/v1"
	model := "gpt-4o-mini"

	systemContent := `You are an expert in Star Trek.`
	userContent := `Who is Jean-Luc Picard?`

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
	}

	_, err = completion.ChatStream(openAIURL, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		}, provider.OpenAI, os.Getenv("OPENAI_API_KEY"))

	if err != nil {
		log.Fatal("😡:", err)
	}

}
