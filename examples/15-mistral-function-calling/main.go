/*
Topic: Parakeet
Generate a chat completion using function calling
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
	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	model := "mistral:7b"

	tools := []llm.Tool{
		{
			Type: "function",
			Function: llm.Function{
				Name:        "hello",
				Description: "Say hello to a given person with his name",
				Parameters: llm.Parameters{
					Type: "object",
					Properties: map[string]llm.Property{
						"name": {
							Type:        "string",
							Description: "The name of the person",
						},
					},
					Required: []string{"name"},
				},
			},
		},
		{
			Type: "function",
			Function: llm.Function{
				Name:        "addNumbers",
				Description: "Make an addition of the two given numbers",
				Parameters: llm.Parameters{
					Type: "object",
					Properties: map[string]llm.Property{
						"a": {
							Type:        "number",
							Description: "first operand",
						},
						"b": {
							Type:        "number",
							Description: "second operand",
						},
					},
					Required: []string{"a", "b"},
				},
			},
		},
	}

	toolsContent, err := llm.GenerateToolsContent(tools)
	if err != nil {
		log.Fatal("😡:", err)
	}

	options := llm.Options{
		Temperature:   0.0,
		RepeatLastN:   2,
		RepeatPenalty: 2.0,
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: toolsContent},
			{Role: "user", Content: llm.GenerateToolsInstruction(`say "hello" to Bob`)},
		},
		Options: options,
		Format:  "json",
		Raw:     true, 
	}

	answer, err := completion.Chat(ollamaUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}
	result, err := llm.PrettyString(answer.Message.Content)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(result)

	query = llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: toolsContent},
			{Role: "user", Content: llm.GenerateToolsInstruction(`add 5 and 40`)},
		},
		Options: options,
		Format:  "json",
		Raw:     true,
	}

	answer, err = completion.Chat(ollamaUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}
	result, err = llm.PrettyString(answer.Message.Content)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(result)
}
