package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/enums/provider"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
)

func main() {

	modelRunnerURL := "http://localhost:12434/engines/llama.cpp/v1"
	model := "ai/smollm2"

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
		option.Temperature:   0.0,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 2.0,
		//option.Verbose:       true,
	})

	query := llm.Query{
		Model:    model,
		Messages: messages,
		Tools:    toolsList,
		Options:  options,
	}

	answer, err := completion.Chat(modelRunnerURL, query, provider.DockerModelRunner)
	if err != nil {
		log.Fatal("😡 completion:", err)
	}
	fmt.Println(answer.Message.ToolCalls)

	// Search tool to call in the answer
	tool, err := answer.Message.ToolCalls.Find("hello")
	if err != nil {
		log.Fatal("😡 ToolCalls.Find:", err)
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
	}

	answer, err = completion.Chat(modelRunnerURL, query, provider.DockerModelRunner)
	if err != nil {
		log.Fatal("😡 completion bis:", err)
	}

	// Search tool to call in the answer
	tool, err = answer.Message.ToolCalls.Find("addNumbers")
	if err != nil {
		log.Fatal("😡 ToolCalls.Find bis:", err)
	}
	result, _ = tool.Function.ToJSONString()
	fmt.Println(result)

}
