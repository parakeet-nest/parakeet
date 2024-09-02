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
	"encoding/json"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
)

/*
This code snippet is a Go program that demonstrates the usage of a chat completion API. It defines a main function that performs the following steps:

1. Sets the ollamaUrl variable to the URL of the chat completion API.
2.Sets the model variable to the name of the model to be used for chat completion.
3. Defines a list of toolsList that contains information about two functions: hello and addNumbers. Each function has a name, description, and parameters.
4. Calls the tools.GenerateContent function to generate the content for the prompt based on the toolsList.
5. Defines an options struct that contains various options for chat completion, such as temperature, repeat last N messages, and repeat penalty.
6. Defines a query struct that specifies the model, messages, options, and format for the chat completion.
7. Calls the completion.Chat function to perform the chat completion using the provided query.
8. Prints the result of the chat completion.
9. Repeats steps 7 and 8 for a different chat completion query.

Overall, this code snippet demonstrates how to use a chat completion API to generate responses based on a given prompt and a list of functions.
*/

func main() {
	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	model := "mistral:7b"

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

	messages := []llm.Message{
		{Role: "user", Content: `say "hello" to Bob`},
	}

	options := llm.Options{
		Temperature:   0.0,
		RepeatLastN:   2,
		RepeatPenalty: 2.0,
	}

	query := llm.Query{
		Model:    model,
		Messages: messages,
		Tools:    toolsList,
		Options:  options,
		Format:   "json",
	}

	answer, err := completion.Chat(ollamaUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}

	// []map[string]interface{}
	toolCalls := answer.Message.ToolCalls
	fmt.Println("👋", toolCalls)


	// Access the first item
	firstItem := toolCalls[0]

	// Access the "function" map
	functionMap, _ := firstItem["function"].(map[string]interface{})
	// Access the "arguments" map
	argumentsMap, _ := functionMap["arguments"].(map[string]interface{})
	// Access the "name" field in arguments
	name, _ := argumentsMap["name"].(string)
	fmt.Println("Name:", name)
	// Access the "name" field in function
	functionName, _ := functionMap["name"].(string)
	fmt.Println("Function name:", functionName)


	jsonBytes, err := json.Marshal(toolCalls)
	if err != nil {
		log.Fatal("😡:", err)
	}
	// Convert JSON bytes to string
	result := string(jsonBytes)

	fmt.Println(result)

	messages = []llm.Message{
		{Role: "user", Content: `add 2 and 40`},
	}

	query = llm.Query{
		Model:    model,
		Messages: messages,
		Tools:    toolsList,
		Options:  options,
		Format:   "json",
	}

	answer, err = completion.Chat(ollamaUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}

	result, err = answer.Message.ToolCallsToJSONString()
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(result)

	// https://github.com/ollama/ollama/blob/main/docs/api.md#response-13
	result, err = answer.Message.FirstToolCallToJSONString()
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(result)
}
