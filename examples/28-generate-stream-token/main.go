/*
Topic: Parakeet
Generate a simple completion with Ollama and parakeet
The output is streamed
*/

package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
	"fmt"
	"log"
)

func main() {
	ollamaUrl := "https://ollamak33g.eu.loclx.io"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	model := "tinydolphin"


	options := llm.Options{
		Temperature: 0.5, // default (0.8)
	}

	firstQuestion := llm.Query{
		Model: model,
		Prompt: "Who is James T Kirk?",
		Options: options,
		TokenHeaderName: "X-TOKEN",
		TokenHeaderValue: "john_doe",
	}
	
	fmt.Println("✋ First Completion:")
	answer, err := completion.GenerateStream(ollamaUrl, firstQuestion,
		func(answer llm.Answer) error {
			fmt.Print(answer.Response)
			return nil
		})

	if err != nil {
		log.Fatal("1 😡:", err)
	}

	secondQuestion := llm.Query{
		Model: model,
		Prompt: "Who is his best friend?",
		Context: answer.Context,
		Options: options,
		TokenHeaderName: "X-TOKEN",
		TokenHeaderValue: "john_doe",
	}

	fmt.Println()
	fmt.Println()
	fmt.Println("✋ Second Completion:")

	_, err = completion.GenerateStream(ollamaUrl, secondQuestion,
		func(answer llm.Answer) error {
			fmt.Print(answer.Response)
			return nil
		})

	if err != nil {
		log.Fatal("2 😡:", err)
	}
}
