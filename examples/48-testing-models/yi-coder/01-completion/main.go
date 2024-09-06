package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
)

/*
https://github.com/01-ai/Yi-Coder/blob/main/cookbook/System_prompt/System_prompt.ipynb
*/

// Here, we will use the example of "writing a quick sort algorithm" to illustrate system prompts.

func main() {
	ollamaUrl := "http://localhost:11434"
	model := "yi-coder:1.5b"

	systemContent := `You are Yi-Coder, you are exceptionally skilled in programming, coding, and any computer-related issues.`

	userContent := `Write a quick sort algorithm in golang.`

	options := llm.Options{
		Temperature: 0.0,
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
	}

	_, err := completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})
	
	if err != nil {
		log.Fatal("😡:", err)
	}


}
