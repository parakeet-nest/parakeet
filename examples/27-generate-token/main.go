/*
Topic: Parakeet
Generate a simple completion with Ollama and parakeet
no streaming
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
		Stop:        []string{},
	}

	firstQuestion := llm.GenQuery{
		Model: model,
		Prompt: "Who is James T Kirk?",
		Options: options,
		TokenHeaderName: "X-TOKEN",
		TokenHeaderValue: "john doe",
	}

	answer, err := completion.Generate(ollamaUrl, firstQuestion)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Response)

	fmt.Println()

	secondQuestion := llm.GenQuery{
		Model: model,
		Prompt: "Who is his best friend?",
		Context: answer.Context,
		Options: options,
		TokenHeaderName: "X-TOKEN",
		TokenHeaderValue: "john doe",
	}

	answer, err = completion.Generate(ollamaUrl, secondQuestion)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Response)


}
