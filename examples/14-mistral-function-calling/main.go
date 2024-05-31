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


	var toolsContent = `[AVAILABLE_TOOLS] [
	{
		"type": "function", 
		"function": {
			"name": "hello",
			"description": "Say hello to a given person with his name",
			"parameters": {
				"type": "object", 
				"properties": {
					"name": {
						"type": "string", 
						"description": "The name of the person"
					}
				}, 
				"required": ["name"]
			}
		}
	},
	{
		"type": "function", 
		"function": {
			"name": "addNumbers",
			"description": "Make an addition of the two given numbers",
			"parameters": {
				"type": "object", 
				"properties": {
					"a": {
						"type": "number", 
						"description": "first operand"
					},
					"b": {
						"type": "number",
						"description": "second operand"
					}
				}, 
				"required": ["a", "b"]
			}
		}
	}
	] [/AVAILABLE_TOOLS]`
	
	options := llm.Options{
		Temperature: 0.0,
		RepeatLastN: 2, 
		RepeatPenalty: 2.0, 
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: toolsContent},
			{Role: "user", Content: `[INST] say "hello" to Bob [/INST]`},
		},
		Options: options,
		Format: "json",
		Raw: true, // override the template
	}
	
	
	answer, err := completion.Chat(ollamaUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Message.Content)
	
	query = llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: toolsContent},
			{Role: "user", Content: `[INST] add 5 and 40 [/INST]`},
		},
		Options: options,
		Format: "json",
		Raw: true,
	}

	answer, err = completion.Chat(ollamaUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Message.Content)
	
}
