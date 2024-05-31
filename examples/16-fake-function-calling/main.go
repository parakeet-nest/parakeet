/*
Topic: Parakeet
Generate a chat completion using fa kind ofunction calling with LLM that does not implement function calling
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
	//model := "mistral:7b"
	model := "phi3:mini"

	systemContentIntroduction := `You have access to the following tools:`

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

	systemContentInstructions := `If the question of the user matched the description of a tool, the tool will be called.
	To call a tool, respond with a JSON object with the following structure: 
	{
	  "name": <name of the called tool>,
	  "arguments": {
	    <name of the argument>: <value of the argument>
	  }
	}
	
	search the name of the tool in the list of tools with the Name field
	`
	
	options := llm.Options{
		Temperature: 0.0,
		RepeatLastN: 2, 
		RepeatPenalty: 2.0, 
		Seed: 123,
		//Stop:        []string{},
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContentIntroduction},
			{Role: "system", Content: toolsContent},
			{Role: "system", Content: systemContentInstructions},
			{Role: "user", Content: `say "hello" to Bob`},
		},
		Options: options,
		Format: "json", // does it works correctly with raw == true?
	}

	answer, err := completion.Chat(ollamaUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Message.Content)

	query = llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContentIntroduction},
			{Role: "system", Content: toolsContent},
			{Role: "system", Content: systemContentInstructions},
			{Role: "user", Content: `add 5 and 40`},
		},
		Options: options,
		Format: "json",
	}

	answer, err = completion.Chat(ollamaUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Message.Content)
}
