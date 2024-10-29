package main

import (
	"os"

	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/llm"

	"log"

	"github.com/parakeet-nest/parakeet/flock"
	"github.com/parakeet-nest/parakeet/ui"
	"github.com/parakeet-nest/parakeet/ui/colors"
)

func PrintMessages(messages []llm.Message) {
	for _, message := range messages {
		ui.Println(colors.Blue, "-", message.Role, ": ", message.Content)
	}
}

func main() {

	ollamaUrl := os.Getenv("OLLAMA_URL")
	if ollamaUrl == "" {
		ollamaUrl = "http://localhost:11434"
	}

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.0,
	})

	tools := []llm.Tool{
		{
			Type: "function",
			Function: llm.Function{
				Name:        "multiplyNumbers",
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

	bob := flock.Agent{
		Name:      "Bob",
		Model:     "allenporter/xlam:1b",
		OllamaUrl: ollamaUrl,
		Options:   options,
		Functions: map[string]flock.AgentFunction{
			"addNumbers": func(args interface{}) (interface{}, error) {
				// Convert the arguments to a map
				argsMap := args.(map[string]interface{})
				a := argsMap["a"].(float64)
				b := argsMap["b"].(float64)
				return a + b, nil
			},
			"multiplyNumbers": func(args interface{}) (interface{}, error) {
				// Convert the arguments to a map
				argsMap := args.(map[string]interface{})
				a := argsMap["a"].(float64)
				b := argsMap["b"].(float64)
				return a * b, nil
			},
		},
	}

	bob.SetInstructions("")

	orchestrator := flock.Orchestrator{}

	bobResponse, err := orchestrator.RunWithTools(
		bob,
		[]llm.Message{
			{Role: "user", Content: `add 2 and 40`},
			{Role: "user", Content: `multiply 2 and 21`},
			{Role: "user", Content: `multiply 34 and 0.5`},
		},
		map[string]interface{}{},
		tools,
		true,
	)

	if err != nil {
		log.Fatalln("😡:", err)
	}

	ui.Println(colors.Green, "🤖 Bob's response: ")
	ui.Println(colors.Yellow, bobResponse.Agent)
	ui.Println(colors.Blue, bobResponse.ContextVariables)
	ui.Println(colors.Purple, bobResponse.Messages)

	PrintMessages(bobResponse.Messages)

	lastBobMessage := bobResponse.Messages[len(bobResponse.Messages)-1]
	ui.Println(colors.Green, "🤖 Last Bob's message: ", lastBobMessage.ToolCalls)

	for _, toolCall := range lastBobMessage.ToolCalls {
		if toolCall.Error == nil {
			ui.Println(colors.Magenta, "🤖", toolCall.Function.Name, "call with ", toolCall.Function.Arguments, "=", toolCall.Result)
		} else {
			ui.Println(colors.Red, "😡", toolCall.Function.Name, "call with ", toolCall.Function.Arguments, "failed with error: ", toolCall.Error)
		}
	}

}
