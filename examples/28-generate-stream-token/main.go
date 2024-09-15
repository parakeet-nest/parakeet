/*
Topic: Parakeet
Generate a simple completion with Ollama and parakeet
The output is streamed
*/

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"

	"github.com/parakeet-nest/parakeet/enums/option"

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

	fmt.Println("✋ First Completion:")
	answer, err := completion.GenerateStream(ollamaUrl, firstQuestion,
		func(answer llm.GenAnswer) error {
			fmt.Print(answer.Response)
			return nil
		})

	if err != nil {
		log.Fatal("1 😡:", err)
	}

	secondQuestion := llm.GenQuery{
		Model:            model,
		Prompt:           "Who is his best friend?",
		Context:          answer.Context,
		Options:          options,
		TokenHeaderName:  "X-TOKEN",
		TokenHeaderValue: os.Getenv("TOKEN"),
	}

	fmt.Println()
	fmt.Println()
	fmt.Println("✋ Second Completion:")

	_, err = completion.GenerateStream(ollamaUrl, secondQuestion,
		func(answer llm.GenAnswer) error {
			fmt.Print(answer.Response)
			return nil
		})

	if err != nil {
		log.Fatal("2 😡:", err)
	}
}
