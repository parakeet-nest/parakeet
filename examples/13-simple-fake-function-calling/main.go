/*
Topic: Parakeet
Generate a chat completion using a kind of function calling
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
	//model := "qwen2:0.5b"
	//model := "qwen2:1.5b"

	model := "phi3:mini"

	var toolsContent = `You have access to the following tools:
	BEGIN LIST
	Name: hello,
	Description: Say hello to a given person with his name
	Parameters: value of name

	Name: addNumbers,
	Description: Make an addition of the two given numbers,
	Parameters: [a, b]
	END LIST
	
	If the question of the user matched the description of a tool, the tool will be called.
	
	To call a tool, respond with a JSON object with the following structure: 
	{
	  "tool": <name of the called tool>,
	  "parameters": <parameters for the tool matching the above Parameters list>
	}
	
	search the name of the tool in the list of tools with the Name field
	`
	
	var systemContent = `You are a helpful AI assistant. The user will enter a sentence.
	If the sentence is near the description of a tool, the assistant will call the tool.
	Output the results in JSON format and trim the spaces of the sentence.`
	
	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.0,
		option.RepeatLastN: 2,
		option.RepeatPenalty: 2.0,
	})
	
	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: toolsContent},
			{Role: "system", Content: systemContent},
			{Role: "user", Content: `add 5 and 40`},
		},
		Options: options,
		Format: "json",
		//Raw: true, // Why?
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
			{Role: "system", Content: systemContent},
			{Role: "user", Content: `say "hello" to Bob`},
		},
		Options: options,
		Format: "json",
		Raw: true, // Why?
	}

	answer, err = completion.Chat(ollamaUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Message.Content)
}
