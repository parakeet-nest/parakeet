package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
)



func main() {
	ollamaUrl := "http://localhost:11434"
	model := "mathstral"

	systemContent := `You are Bob, you are exceptionally skilled in mathematics.`

	userContent := `solve : 2x-5=3x+2`

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
