/*
Topic: Parakeet
Generate a chat completion with Ollama and parakeet
no streaming
*/

package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	//"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/enums/provider"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
)

func main() {
	modelRunnerURL := "http://localhost:12434/engines/llama.cpp/v1"
	model := "ai/qwen2.5:latest" 

	systemContent := `You are an expert in Star Trek.`
	userContent := `Who is Jean-Luc Picard?`

	/*
	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.5,
		option.RepeatPenalty: 2.0,
	})
	*/

	options := llm.Options{
		Temperature: 	0.5,
		RepeatPenalty: 	2.0,
	}

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
