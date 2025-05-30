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
	//ollamaUrl := "http://robby.local:4000" // this is my RPI5
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	model := "allenporter/xlam:1b"

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

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.0,
		option.RepeatLastN: 2,
		option.RepeatPenalty: 2.0,
	})


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

	// Search tool to call in the answer
	tool, err := answer.Message.ToolCalls.Find("hello")
	if err != nil {
		log.Fatal("😡:", err)
	}
	result, _ := tool.Function.ToJSONString()
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

	// Search tool to call in the answer
	tool, err = answer.Message.ToolCalls.Find("addNumbers")
	if err != nil {
		log.Fatal("😡:", err)
	}
	result, _ = tool.Function.ToJSONString()
	fmt.Println(result)

}
