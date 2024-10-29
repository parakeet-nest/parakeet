package main

import (
	"fmt"

	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/flock"
	"github.com/parakeet-nest/parakeet/llm"
)

func main() {
	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.0,
	})

	// Create agent with function capability
	calculator := flock.Agent{
		Name:    "Calculator",
		Model:   "allenporter/xlam:1b",
		OllamaUrl: "http://localhost:11434",
		Options: options,
		Functions: map[string]flock.AgentFunction{
			"multiply": func(args interface{}) (interface{}, error) {
				argsMap := args.(map[string]interface{})
				return argsMap["a"].(float64) * argsMap["b"].(float64), nil
			},
		},
	}

	// Define tool
	tools := []llm.Tool{
		{
			Type: "function",
			Function: llm.Function{
				Name:        "multiply",
				Description: "Make a multiplication of the two given numbers",
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

	// Run calculation
	orchestrator := flock.Orchestrator{}

	response, _ := orchestrator.RunWithTools(
		calculator,
		[]llm.Message{{Role: "user", Content: "multiply 5 and 3"}},
		map[string]interface{}{},
		tools,
		true,
	)

	// Access result
	fmt.Println("Result:", response.GetLastMessage().ToolCalls[0].Result)

}
