/*
Topic: Parakeet
Generate a chat completion with Docker Model Runner and parakeet
The output is streamed
*/

package main

import (
	"fmt"
	"log"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/enums/provider"
	"github.com/parakeet-nest/parakeet/llm"
)



func main() {

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.5,
		option.RepeatPenalty: 3.0,
	})

	query := llm.Query{
		Model: "ai/qwen2.5:latest",
		Messages: []llm.Message{
			{Role: "system", Content: `You are a Borg in Star Trek. Speak like a Borg`},
			{Role: "user", Content: `Who is Jean-Luc Picard?`},
		},
		Options: options,
	}

	_, err := completion.ChatStream("http://localhost:12434/engines/llama.cpp/v1", query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		}, provider.DockerModelRunner)

	if err != nil {
		log.Fatal("😡:", err)
	}
}
