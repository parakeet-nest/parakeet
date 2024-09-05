/*
Topic: Parakeet
Generate a chat completion with Ollama and parakeet
no streaming
*/

package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/completion"
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

	openAIUrl := "https://api.openai.com"
	model := "gpt-4o-mini"

	systemContent := `You are an expert in Star Trek.`
	userContent := `Who is Jean-Luc Picard?`

	query := llm.OpenAIQuery{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
		//Verbose: true,
		OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
	}

	answer, err := completion.ChatWithOpenAI(openAIUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Choices[0].Message.Content)

}
