/*
Topic: Parakeet
Generate a chat completion with Ollama and parakeet
no streaming
*/

package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/enums/option"


	"fmt"
	"log"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	//model := "deepseek-coder"
	model := "qwen2.5-coder:7b" 

	systemContent := `You are an expert in computer programming.
	Please make friendly answer for the noobs.
	Add source code examples if you can.`

	userContent := `I need a clear explanation regarding the following question:
	Can you create a "hello world" program in Golang?
	And, please, be structured with bullet points`


	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.5,
		option.RepeatLastN: 2,
		option.RepeatPenalty: 2.0,
		option.Verbose: true,
	})

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
	}

	answer, err := completion.Chat(ollamaUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Message.Content)

	//fmt.Println("[llm/query]", query.ToJsonString())
	//fmt.Println("[llm/completion]", answer.ToJsonString())

}
