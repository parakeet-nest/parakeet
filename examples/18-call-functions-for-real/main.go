/*
Topic: Parakeet
Generate a chat completion using function calling
no streaming
This example:
- uses Mistral model
- make a list of tools
- use the list of tools to generate content for the prompt
- retrieve the function from the list of tools
*/

package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/gear"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/tools"
	"github.com/parakeet-nest/parakeet/wasm"

	"github.com/parakeet-nest/parakeet/enums/option"


	"fmt"
	"log"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	//ollamaUrl := "http://bob.local:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	//model := "mistral:7b"
	model := "mistral:latest"

	wasmPlugin, _ := wasm.NewPlugin("./wasm/plugin.wasm", nil)

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

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.0,
		option.RepeatLastN: 2,
		option.RepeatPenalty: 2.0,
	})

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: toolsContent},
			{Role: "user", Content: tools.GenerateUserToolsInstructions(`say "hello" to Bob`)},
		},
		Options: options,
		Format:  "json",
		Raw:     true,
	}

	answer, err := completion.Chat(ollamaUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}

	jsonRes, err := gear.JSONParse(answer.Message.Content)

	if err != nil {
		log.Fatal("😡:", err)
	}
	functionName := jsonRes["name"].(string)
	name := jsonRes["arguments"].(map[string]interface{})["name"].(string)

	fmt.Println("Calling", functionName, "with", name)

	// call the function of the wasm plugin
	res, _ := wasmPlugin.Call(functionName, []byte(name))

	// display the result
	fmt.Println(string(res))

	query = llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "system", Content: toolsContent},
			{Role: "user", Content: tools.GenerateUserToolsInstructions(`add 2 and 40`)},
		},
		Options: options,
		Format:  "json",
		Raw:     true,
	}

	answer, err = completion.Chat(ollamaUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}

	jsonRes, err = gear.JSONParse(answer.Message.Content)
	if err != nil {
		log.Fatal("😡:", err)
	}

	functionName = jsonRes["name"].(string)
	arguments := jsonRes["arguments"].(map[string]interface{})
	fmt.Println("Calling", functionName, "with", arguments)

	// call the function of the wasm plugin
	res, _ = wasmPlugin.Call(functionName, []byte(gear.JSONStringify(arguments)))

	// display the result
	fmt.Println(string(res))

}
