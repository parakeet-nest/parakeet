/*
Topic: Parakeet
Generate a chat completion with Ollama and parakeet
no streaming
*/

package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/enums/provider"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
)

func main() {
	modelRunnerURL := "http://localhost:12434/engines/llama.cpp/v1"

	model := "ai/qwen2.5:latest" 

	systemContent := `You are an expert in computer programming.
	Please make friendly answer for the noobs.
	Add source code examples if you can.`

	userContent := `I need a clear explanation regarding the following question:
	Can you create a "hello world" program in Golang?
	And, please, be structured with bullet points`


	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.5,
		option.RepeatPenalty: 2.0,
	})

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
	}

	answer, err := completion.Chat(modelRunnerURL, query, provider.DockerModelRunner)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Message.Content)

}
