/*
Topic: Parakeet
Generate a simple completion with Ollama and parakeet
no streaming
*/
package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/enums/option"


	"fmt"
	"log"
)

func main() {
	// create a `.env` file with the following content:
	// TOKEN=your_token
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡:", err)
	}

	ollamaUrl := "https://ollama.wasm.ninja"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	model := "tinydolphin"


	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.5,
	})

	firstQuestion := llm.GenQuery{
		Model:            model,
		Prompt:           "Who is James T Kirk?",
		Options:          options,
		TokenHeaderName:  "X-TOKEN",
		TokenHeaderValue: os.Getenv("TOKEN"),
	}

	answer, err := completion.Generate(ollamaUrl, firstQuestion)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Response)

	fmt.Println()

	secondQuestion := llm.GenQuery{
		Model:            model,
		Prompt:           "Who is his best friend?",
		Context:          answer.Context,
		Options:          options,
		TokenHeaderName:  "X-TOKEN",
		TokenHeaderValue: os.Getenv("TOKEN"),
	}

	answer, err = completion.Generate(ollamaUrl, secondQuestion)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Response)

}
