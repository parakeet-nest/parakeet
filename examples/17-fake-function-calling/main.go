/*
Topic: Parakeet
Generate a chat completion using fa kind ofunction calling with LLM that does not implement function calling
no streaming
*/

package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/tools"
	"github.com/parakeet-nest/parakeet/gear"
	"github.com/parakeet-nest/parakeet/enums/option"


	"fmt"
	"log"
)

func main() {
	//ollamaUrl := "http://localhost:11434"
	//ollamaUrl := "http://bob.local:11434"
	// if working from a container
	ollamaUrl := "http://host.docker.internal:11434"
	//model := "mistral:7b"
	//model := "phi3:mini"
	model := "qwen2.5:0.5b"

	systemContentIntroduction := `You have access to the following tools:`

	toolsList := []llm.Tool{
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

	toolsContent, err := tools.GenerateAvailableToolsContent(toolsList)
	if err != nil {
		log.Fatal("😡:", err)
	}

	/*
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
	*/
	
	// ✋ Use it only if the LLM does not implement function calling.
	systemContentInstructions := tools.GenerateSystemToolsInstructions()

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.0,
		option.RepeatLastN: 2,
		option.RepeatPenalty: 2.0,
		option.Seed: 123,
	})

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContentIntroduction},
			{Role: "system", Content: toolsContent},
			{Role: "system", Content: systemContentInstructions},
			{Role: "user", Content: `say "hello" to Bob`},
		},
		Options: options,
		Format:  "json",
	}

	answer, err := completion.Chat(ollamaUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}

	result, err := gear.PrettyString(answer.Message.Content)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(result)

	query = llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: systemContentIntroduction},
			{Role: "system", Content: toolsContent},
			{Role: "system", Content: systemContentInstructions},
			{Role: "user", Content: `add 5 and 40`},
		},
		Options: options,
		Format:  "json",
	}

	answer, err = completion.Chat(ollamaUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}

	result, err = gear.PrettyString(answer.Message.Content)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(result)

}
